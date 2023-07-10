package testcase

import (
	"fmt"

	"github.com/rancher/rke2/tests/acceptance/core/service/factory"
	"github.com/rancher/rke2/tests/acceptance/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestBuildCluster test the creation of a cluster using terraform
func TestBuildCluster(g GinkgoTInterface, destroy bool) {
	cluster := factory.GetCluster(g)

	Expect(cluster.Status).To(Equal("cluster created"))
	Expect(shared.KubeConfigFile).ShouldNot(BeEmpty())
	Expect(cluster.ServerIPs).ShouldNot(BeEmpty())

	fmt.Println("\nKubeconfig file:")
	shared.PrintFileContents(shared.KubeConfigFile)
	fmt.Println("Base64 Encoded Kubeconfig file:")
	shared.PrintBase64Encoded(shared.KubeConfigFile)
	fmt.Println(
		"\nServer Node IPS:", cluster.ServerIPs,
		"\nAgent Node IPS:", cluster.AgentIPs,
		"\nWindows Agent Node IPS:", cluster.WinAgentIPs,
	)

	if cluster.NumAgents > 0 {
		Expect(cluster.AgentIPs).ShouldNot(BeEmpty())
	} else {
		Expect(cluster.AgentIPs).Should(BeEmpty())
	}

	if cluster.NumWinAgents > 0 {
		Expect(cluster.WinAgentIPs).ShouldNot(BeEmpty())
	} else {
		Expect(cluster.WinAgentIPs).Should(BeEmpty())
	}
}

// TestSonobuoyMixedOS runs sonobuoy tests for mixed os cluster (linux + windows) node
func TestSonobuoyMixedOS(version string, delete bool) {
	err := shared.InstallSonobuoyMixedOS(version)
	if err != nil {
		fmt.Errorf("Error installing sonobuoy: ", err)
		return
	}
	
	cmd := "sonobuoy run --kubeconfig=" + shared.KubeConfigFile +
		" --plugin my-sonobuoy-plugins/mixed-workload-e2e/mixed-workload-e2e.yaml" + 
		" --aggregator-node-selector kubernetes.io/os:linux --wait"
	res, err := shared.RunCommandHost(cmd)
	Expect(err).NotTo(HaveOccurred(), "failed output: " + res)
	
	cmd = fmt.Sprintf("sonobuoy retrieve --kubeconfig=%s",shared.KubeConfigFile)
	testResultTar, err := shared.RunCommandHost(cmd)
	Expect(err).NotTo(HaveOccurred(), "failed cmd: "+ cmd)
	
	cmd = fmt.Sprintf("sonobuoy results %s",testResultTar)
	res, err = shared.RunCommandHost(cmd)
	Expect(err).NotTo(HaveOccurred(), "failed cmd: "+ cmd)
	Expect(res).Should(ContainSubstring("Plugin: mixed-workload-e2e\nStatus: passed\n"))

	if delete{
		cmd = fmt.Sprintf("sonobuoy delete --all --wait --kubeconfig=%s", shared.KubeConfigFile)
		res, err = shared.RunCommandHost(cmd)
		Expect(err).NotTo(HaveOccurred(), "failed cmd: "+ cmd)
	}
}
