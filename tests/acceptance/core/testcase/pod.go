package testcase

import (
	"fmt"
	"strings"

	"github.com/rancher/rke2/tests/acceptance/core/service/assert"
	"github.com/rancher/rke2/tests/acceptance/shared"

	. "github.com/onsi/gomega"
)

const statusCompleted = "Completed"

// TestPodStatus test the status of the pods in the cluster using 2 custom assert functions
func TestPodStatus(
	podAssertRestarts assert.PodAssertFunc,
	podAssertReady assert.PodAssertFunc,
	podAssertStatus assert.PodAssertFunc,
) {
	fmt.Printf("\nChecking pod status")
	Eventually(func(g Gomega) {
		pods, err := shared.ParsePods(false)
		g.Expect(err).NotTo(HaveOccurred())

		for _, pod := range pods {
			fmt.Printf(".")
			if strings.Contains(pod.Name, "helm-install") {
				g.Expect(pod.Status).Should(Equal(statusCompleted), pod.Name)
			} else if strings.Contains(pod.Name, "apply") &&
				strings.Contains(pod.NameSpace, "system-upgrade") {
				g.Expect(pod.Status).Should(SatisfyAny(
					ContainSubstring("Error"),
					Equal(statusCompleted),
				), pod.Name)
			} else {
				g.Expect(pod.Status).Should(Equal("Running"), pod.Name)
				if podAssertRestarts != nil {
					podAssertRestarts(g, pod)
				}
				if podAssertReady != nil {
					podAssertReady(g, pod)
				}
				if podAssertStatus != nil {
					podAssertStatus(g, pod)
				}
			}
		}
	}, "600s", "3s").Should(Succeed())

	fmt.Println("\nCluster pods: ")
	_, err := shared.ParsePods(true)
	if err != nil {
		fmt.Println("Error retrieving pods: ", err)
	}
}
