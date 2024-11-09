package tests

import (
	"Transactio/internal/fileStorage/utils"
	"reflect"
	"testing"
)

func TestCypher(t *testing.T) {

	testData := "Test data"
	testPassword := "testPassword"

	type args struct {
		data     []byte
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "#1",
			args: args{
				data:     []byte(testData),
				password: testPassword,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := tt.args.password
			encData, _ := utils.EncryptData(tt.args.data, pass)
			decrData, _ := utils.DecryptData(encData, pass)
			if !reflect.DeepEqual(tt.args.data, decrData) {
				t.Errorf("Test failed. Data isn't same. \nPass = %v, \nData = %v, \nEncData = %v, \nDecrData = %v.",
					tt.args.password, string(tt.args.data), string(encData), string(decrData))
				return
			}

			t.Logf("Test comlited! \nPass = %v, \nData = %v, \nEncData = %v, \nDecrData = %v.",
				tt.args.password, string(tt.args.data), string(encData), string(decrData))
		})
	}

}
