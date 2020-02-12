package provider

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/config"

	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	metricAnnotationPrefix = "wavefront.com.external.metric"
)

type hpaListener struct {
	kubeClient kubernetes.Interface
	addFunc    RuleHandlerFunc
	deleteFunc RuleHandlerFunc
}

func StartHPAListener(client kubernetes.Interface, addFunc, deleteFunc RuleHandlerFunc) {
	listener := &hpaListener{
		kubeClient: client,
		addFunc:    addFunc,
		deleteFunc: deleteFunc,
	}
	go listener.listen()
}

func (l *hpaListener) listen() {
	log.Info("listening for HPA instances")

	rc := l.kubeClient.AutoscalingV2beta1().RESTClient()
	lw := cache.NewListWatchFromClient(rc, "horizontalpodautoscalers", v1.NamespaceAll, fields.Everything())
	inf := cache.NewSharedInformer(lw, &v2beta1.HorizontalPodAutoscaler{}, 0)

	inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			hpa := obj.(*v2beta1.HorizontalPodAutoscaler)
			rules := rulesFromAnnotations(hpa.Annotations)
			if len(rules) > 0 {
				l.addFunc(rules)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldHPA := oldObj.(*v2beta1.HorizontalPodAutoscaler)
			newHPA := newObj.(*v2beta1.HorizontalPodAutoscaler)

			// HPA objects are updated frequently when status changes
			// validate if annotations have changed
			if reflect.DeepEqual(oldHPA.Annotations, newHPA.Annotations) {
				log.Debugf("annotations have not changed for %s", newHPA.Name)
				return
			}

			oldRules := rulesFromAnnotations(oldHPA.Annotations)
			if len(oldRules) > 0 {
				l.deleteFunc(oldRules)
			}
			newRules := rulesFromAnnotations(newHPA.Annotations)
			if len(newRules) > 0 {
				l.addFunc(newRules)
			}
		},
		DeleteFunc: func(obj interface{}) {
			hpa := obj.(*v2beta1.HorizontalPodAutoscaler)
			rules := rulesFromAnnotations(hpa.Annotations)
			if len(rules) > 0 {
				l.deleteFunc(rules)
			}
		},
	})
	go inf.Run(wait.NeverStop)
}

func rulesFromAnnotations(annotations map[string]string) []config.MetricRule {
	plen := len(metricAnnotationPrefix)
	var rules []config.MetricRule
	for k, v := range annotations {
		if strings.HasPrefix(k, metricAnnotationPrefix) {
			if len(k) > plen+1 {
				rules = append(rules, config.MetricRule{
					Name:  k[plen+1:],
					Query: v,
				})
			}
		}
	}
	return rules
}
