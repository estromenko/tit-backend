package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/models"
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedNetworkingV1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
)

const dashboardPort = 8888

type DashboardData struct {
	Password string `json:"password"`
}

type DashboardService struct {
	db              *core.Database
	conf            *core.Config
	clientSet       *kubernetes.Clientset
	podsClient      typedCoreV1.PodInterface
	servicesClient  typedCoreV1.ServiceInterface
	ingressesClient typedNetworkingV1.IngressInterface
}

func NewDashboardService(db *core.Database, conf *core.Config) (*DashboardService, error) {
	var (
		clientConfig *rest.Config
		err          error
	)

	if conf.KubernetesUseInClusterConfig {
		clientConfig, err = rest.InClusterConfig()
	} else {
		clientConfig, err = clientcmd.BuildConfigFromFlags("", conf.KubernetesConfigPath)
	}

	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	podsClient := clientSet.CoreV1().Pods(conf.KubernetesDashboardNamespace)
	servicesClient := clientSet.CoreV1().Services(conf.KubernetesDashboardNamespace)
	ingressesClient := clientSet.NetworkingV1().Ingresses(conf.KubernetesDashboardNamespace)

	return &DashboardService{
		db:              db,
		conf:            conf,
		clientSet:       clientSet,
		podsClient:      podsClient,
		servicesClient:  servicesClient,
		ingressesClient: ingressesClient,
	}, nil
}

func (d *DashboardService) createPodForUser(user *models.User) *coreV1.Pod {
	resourceName := fmt.Sprintf("tit-dashboard-%d", user.ID)

	return &coreV1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": resourceName,
			},
		},
		Spec: coreV1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			Containers: []coreV1.Container{
				{
					Name:  "dashboard",
					Image: d.conf.DashboardImage,
					Env: []coreV1.EnvVar{
						{Name: "PASSWORD", Value: user.DashboardPassword},
					},
					Ports: []coreV1.ContainerPort{
						{
							ContainerPort: dashboardPort,
						},
					},
				},
			},
		},
	}
}

func (d *DashboardService) createServiceForUser(user *models.User) *coreV1.Service {
	resourceName := fmt.Sprintf("tit-dashboard-%d", user.ID)

	return &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      resourceName,
			Namespace: coreV1.NamespaceDefault,
			Labels: map[string]string{
				"app": resourceName,
			},
		},
		Spec: coreV1.ServiceSpec{
			Ports: []coreV1.ServicePort{
				{
					Port: dashboardPort,
					TargetPort: intstr.IntOrString{
						IntVal: dashboardPort,
					},
				},
			},
			Selector: map[string]string{
				"app": resourceName,
			},
		},
	}
}

func (d *DashboardService) createIngressForUser(user *models.User) *networkingV1.Ingress {
	resourceName := fmt.Sprintf("tit-dashboard-%d", user.ID)
	pathTypePrefix := networkingV1.PathTypePrefix

	return &networkingV1.Ingress{
		ObjectMeta: metaV1.ObjectMeta{
			Name: resourceName,
		},
		Spec: networkingV1.IngressSpec{
			Rules: []networkingV1.IngressRule{
				{
					Host: d.conf.DashboardIngressDomain,
					IngressRuleValue: networkingV1.IngressRuleValue{
						HTTP: &networkingV1.HTTPIngressRuleValue{
							Paths: []networkingV1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/%d", user.ID),
									PathType: &pathTypePrefix,
									Backend: networkingV1.IngressBackend{
										Service: &networkingV1.IngressServiceBackend{
											Name: resourceName,
											Port: networkingV1.ServiceBackendPort{
												Number: dashboardPort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DashboardService) generateRandomPassword() string {
	return uuid.New().String()
}

func (d *DashboardService) StartDashboard(ctx context.Context, user *models.User) error {
	resourceName := fmt.Sprintf("tit-dashboard-%d", user.ID)

	dashboardPassword := d.generateRandomPassword()
	user.DashboardPassword = dashboardPassword

	pod := d.createPodForUser(user)

	_, err := d.podsClient.Create(ctx, pod, metaV1.CreateOptions{})
	if err != nil {
		return err
	}

	service := d.createServiceForUser(user)

	_, err = d.servicesClient.Create(ctx, service, metaV1.CreateOptions{})
	if err != nil {
		_ = d.podsClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})

		return err
	}

	ingress := d.createIngressForUser(user)

	_, err = d.ingressesClient.Create(ctx, ingress, metaV1.CreateOptions{})
	if err != nil {
		_ = d.podsClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})
		_ = d.servicesClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})

		return err
	}

	_, err = d.db.NewUpdate().
		Model(user).
		Where("id = ?", user.ID).
		Set("dashboard_password = ?", dashboardPassword).
		Exec(ctx)
	if err != nil {
		_ = d.podsClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})
		_ = d.servicesClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})
		_ = d.ingressesClient.Delete(ctx, resourceName, metaV1.DeleteOptions{})

		return err
	}

	return err
}

func (d *DashboardService) IsDashboardRunning(ctx context.Context, user *models.User) bool {
	_, err := d.ingressesClient.Get(ctx, fmt.Sprintf("tit-dashboard-%d", user.ID), metaV1.GetOptions{})

	return err == nil
}
