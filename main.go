package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang/glog"
	"github.com/rs/cors"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfigFile = flag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	port           = flag.String("port", "8080", "Application port to use")
)

func init() {
	flag.Parse()
}

type KubeQLClient struct {
	client *kubernetes.Clientset
}

func main() {
	fmt.Println("Hi from kubeql")

	kubeQLClient := buildKubeQL()

	resolver := &Resolver{kubeQLClient: kubeQLClient}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(NewExecutableSchema(Config{Resolvers: resolver})))

	glog.V(5).Infof("Connect to http://localhost:%s/ for GraphQL playground", *port)
	glog.Fatalf("Fatal: %s", http.ListenAndServe(":"+*port, router))
}

func getKubeConfig() *rest.Config {
	if *kubeConfigFile != "" {
		glog.V(1).Infof("Using kubeconfig file: %s", *kubeConfigFile)
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigFile)
		if err != nil {
			glog.Fatalf("Failed to build config: %v", err)
		}
		return config
	}

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	}

	return kubeConfig
}

func buildKubeQL() *KubeQLClient {
	config := getKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		glog.Fatalf("Failed to build config: %v", err)
	}

	kubeQL := &KubeQLClient{client: clientset}

	return kubeQL
}
