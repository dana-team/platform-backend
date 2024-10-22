package mocks

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const dot = "."

// prepareCNAMERecord returns a mocked CNAME object.
func prepareCNAMERecord(name, cappName, cappNSName, hostname string, readyStatus, syncedStatus corev1.ConditionStatus) dnsrecordv1alpha1.CNAMERecord {
	cnameRecord := prepareBaseCNAMERecord(name, cappName, cappNSName, hostname)
	cnameRecord.Status.ResourceStatus = xpv1.ResourceStatus{ConditionedStatus: xpv1.ConditionedStatus{
		Conditions: []xpv1.Condition{
			{Type: xpv1.TypeSynced, Status: syncedStatus},
			{Type: xpv1.TypeReady, Status: readyStatus},
		},
	},
	}
	return cnameRecord
}

// prepareBaseCNAMERecord returns a mocked CNAME object without conditions.
func prepareBaseCNAMERecord(name, cappName, cappNSName, hostname string) dnsrecordv1alpha1.CNAMERecord {
	hostnameID := hostname + dot

	return dnsrecordv1alpha1.CNAMERecord{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{testutils.ParentCappNSLabel: cappNSName, testutils.ParentCappLabel: cappName},
		},
		Status: dnsrecordv1alpha1.CNAMERecordStatus{
			AtProvider:     dnsrecordv1alpha1.CNAMERecordObservation{ID: &hostnameID},
			ResourceStatus: xpv1.ResourceStatus{},
		},
	}
}
