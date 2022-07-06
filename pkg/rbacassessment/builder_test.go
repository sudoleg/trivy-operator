package rbacassessment_test

import (
	"testing"

	"github.com/aquasecurity/trivy-operator/pkg/apis/aquasecurity/v1alpha1"
	"github.com/aquasecurity/trivy-operator/pkg/configauditreport"
	"github.com/aquasecurity/trivy-operator/pkg/rbacassessment"
	"github.com/aquasecurity/trivy-operator/pkg/trivyoperator"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
)

func TestReportBuilder(t *testing.T) {

	t.Run("Should build report for namespaced resource", func(t *testing.T) {
		g := NewGomegaWithT(t)

		report, err := rbacassessment.NewReportBuilder(scheme.Scheme).
			Controller(&rbacv1.Role{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-owner",
					Namespace: "qa",
				},
				Rules: []rbacv1.PolicyRule{},
			}).
			ResourceSpecHash("xyz").
			PluginConfigHash("nop").
			Data(v1alpha1.RbacAssessmentReportData{}).
			GetReport()
		g.Expect(err).ToNot(HaveOccurred())
		assessmentReport := rbacReport()
		g.Expect(report).To(Equal(assessmentReport))
	})

	t.Run("Should build report for cluster scoped resource", func(t *testing.T) {
		g := NewGomegaWithT(t)

		report, err := configauditreport.NewReportBuilder(scheme.Scheme).
			Controller(&rbacv1.ClusterRole{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "system:controller:node-controller",
				},
			}).
			ResourceSpecHash("xyz").
			PluginConfigHash("nop").
			Data(v1alpha1.ConfigAuditReportData{}).
			GetClusterReport()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(report).To(Equal(v1alpha1.ClusterConfigAuditReport{
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterrole-6f69bb5b79",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         "rbac.authorization.k8s.io/v1",
						Kind:               "ClusterRole",
						Name:               "system:controller:node-controller",
						Controller:         pointer.BoolPtr(true),
						BlockOwnerDeletion: pointer.BoolPtr(false),
					},
				},
				Labels: map[string]string{
					trivyoperator.LabelResourceKind:      "ClusterRole",
					trivyoperator.LabelResourceNameHash:  "6f69bb5b79",
					trivyoperator.LabelResourceNamespace: "",
					trivyoperator.LabelResourceSpecHash:  "xyz",
					trivyoperator.LabelPluginConfigHash:  "nop",
				},
				Annotations: map[string]string{
					trivyoperator.LabelResourceName: "system:controller:node-controller",
				},
			},
			Report: v1alpha1.ConfigAuditReportData{},
		}))
	})
}

func rbacReport() v1alpha1.RbacAssessmentReport {
	return v1alpha1.RbacAssessmentReport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "role-some-owner",
			Namespace: "qa",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "rbac.authorization.k8s.io/v1",
					Kind:               "Role",
					Name:               "some-owner",
					Controller:         pointer.BoolPtr(true),
					BlockOwnerDeletion: pointer.BoolPtr(false),
				},
			},
			Labels: map[string]string{
				trivyoperator.LabelResourceKind:      "Role",
				trivyoperator.LabelResourceName:      "some-owner",
				trivyoperator.LabelResourceNamespace: "qa",
				trivyoperator.LabelResourceSpecHash:  "xyz",
				trivyoperator.LabelPluginConfigHash:  "nop",
			},
		},
		Report: v1alpha1.RbacAssessmentReportData{},
	}
}
