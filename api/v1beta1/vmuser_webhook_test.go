package v1beta1

import (
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestVMUser_sanityCheck(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       VMUserSpec
		Status     VMUserStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid auths",
			fields: fields{
				Spec: VMUserSpec{
					UserName:    pointer.StringPtr("user"),
					BearerToken: pointer.StringPtr("bearer"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ref",
			fields: fields{
				Spec: VMUserSpec{
					UserName: pointer.StringPtr("some-user"),
					TargetRefs: []TargetRef{
						{
							CRD:    &CRDRef{Name: "sm"},
							Static: &StaticRef{URL: "some"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ref wo targets",
			fields: fields{
				Spec: VMUserSpec{
					UserName: pointer.StringPtr("some-user"),
					TargetRefs: []TargetRef{
						{
							Paths: []string{"/some-path"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ref crd, bad kind",
			fields: fields{
				Spec: VMUserSpec{
					UserName: pointer.StringPtr("some-user"),
					TargetRefs: []TargetRef{
						{
							CRD: &CRDRef{
								Name:      "some-1",
								Kind:      "badkind",
								Namespace: "some-ns",
							},
							Paths: []string{"/some-path"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ref crd, bad empty ns",
			fields: fields{
				Spec: VMUserSpec{
					UserName: pointer.StringPtr("some-user"),
					TargetRefs: []TargetRef{
						{
							CRD: &CRDRef{
								Name:      "some-1",
								Kind:      "VMSingle",
								Namespace: "",
							},
							Paths: []string{"/some-path"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "correct crd target",
			fields: fields{
				Spec: VMUserSpec{
					TargetRefs: []TargetRef{
						{
							CRD: &CRDRef{
								Name:      "some-1",
								Namespace: "some-ns",
								Kind:      "VMSingle",
							},
							Paths: []string{"/"},
						},
						{
							Static: &StaticRef{
								URL: "http://some-url",
							},
							Paths: []string{"/targets"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &VMUser{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if err := cr.sanityCheck(); (err != nil) != tt.wantErr {
				t.Errorf("sanityCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
