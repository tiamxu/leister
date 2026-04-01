package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/tiamxu/kit/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Tool struct{}

func (t *Tool) Name() string        { return "kube" }
func (t *Tool) Description() string { return "Manage kubernetes resources" }

func (t *Tool) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag("namespace", "n", "default", "Kubernetes namespace"),
		cli.StringFlag("name", "N", "", "Deployment name"),
		cli.StringFlag("kubeconfig", "k", "~/.kube/config", "Kubernetes config file"),
	}
}

func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("get").
			SetDescription("Get k8s resource deployment").
			AddFlags(cli.RequiredFlag(cli.StringFlag("name", "N", "", "Deployment name"))).
			SetRun(func(ctx *cli.Context) error {
				return RunGetDeployment(ctx)
			}),
		cli.NewCommand("restart").
			SetDescription("Restart k8s resource deployment").
			AddFlags(cli.RequiredFlag(cli.StringFlag("name", "N", "", "Deployment name"))).
			SetRun(func(ctx *cli.Context) error {
				return RunRestart(ctx)
			}),
		cli.NewCommand("create").
			SetDescription("Create resource deployment").
			AddFlags(cli.RequiredFlag(cli.StringFlag("name", "N", "", "Deployment name"))).
			SetRun(func(ctx *cli.Context) error {
				return CreateDeployment(ctx)
			}),
	}
}

func getClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %v", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return client, nil
}

func RunGetDeployment(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	kubeconfig := ctx.String("kubeconfig")

	client, err := getClient(kubeconfig)
	if err != nil {
		return err
	}

	deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	fmt.Printf("Deployment: %s\n", deployment.Name)
	fmt.Printf("Replicas: %d/%d\n", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
	fmt.Printf("Image: %s\n", deployment.Spec.Template.Spec.Containers[0].Image)
	fmt.Printf("Status: %s\n", deployment.Status.Conditions[len(deployment.Status.Conditions)-1].Type)

	return nil
}

func RunRestart(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	kubeconfig := ctx.String("kubeconfig")

	client, err := getClient(kubeconfig)
	if err != nil {
		return err
	}

	// 通过更新 annotations 触发重启
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}
	deployment.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = client.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %v", err)
	}

	fmt.Printf("Deployment %s restarted successfully\n", name)
	return nil
}

func CreateDeployment(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	kubeconfig := ctx.String("kubeconfig")

	client, err := getClient(kubeconfig)
	if err != nil {
		return err
	}

	// 创建一个简单的 deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: "nginx:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = client.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	fmt.Printf("Deployment %s created successfully\n", name)
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
