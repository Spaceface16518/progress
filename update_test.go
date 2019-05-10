package progress

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func generateTestIds(amount int) []uuid.UUID {
	s := make([]uuid.UUID, amount)

	for i := 0; i < amount; i++ {
		s[i] = uuid.Must(uuid.NewV4())
	}
	return s
}

func Test_fixedUpdate_ID(t *testing.T) {
	testIds := generateTestIds(3)

	type fields struct {
		id       uuid.UUID
		progress float32
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			"ID #1",
			fields{
				id:       testIds[0],
				progress: 0,
			},
			testIds[0],
		},
		{
			"ID #2",
			fields{
				id:       testIds[1],
				progress: 0,
			},
			testIds[1],
		},
		{
			"ID #3",
			fields{
				id:       testIds[2],
				progress: 0,
			},
			testIds[2],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &fixedUpdate{
				id:       tt.fields.id,
				progress: tt.fields.progress,
			}
			if got := update.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fixedUpdate.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixedUpdate_NewVal(t *testing.T) {
	testIds := generateTestIds(3)
	type fields struct {
		id       uuid.UUID
		progress float32
	}
	type args struct {
		in0 float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		{
			"To 0",
			fields{
				testIds[0],
				0,
			},
			args{
				1,
			},
			0,
		},
		{
			"To 100",
			fields{
				testIds[1],
				1,
			},
			args{
				0.5,
			},
			1,
		},
		{
			"To 50",
			fields{
				testIds[2],
				0.5,
			},
			args{
				0,
			},
			0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &fixedUpdate{
				id:       tt.fields.id,
				progress: tt.fields.progress,
			}
			if got := update.NewVal(tt.args.in0); got != tt.want {
				t.Errorf("fixedUpdate.NewVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deltaUpdate_ID(t *testing.T) {
	testIds := generateTestIds(3)
	type fields struct {
		id       uuid.UUID
		progress float32
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			"ID #1",
			fields{
				testIds[0],
				0,
			},
			testIds[0],
		},
		{
			"ID #2",
			fields{
				testIds[1],
				0,
			},
			testIds[1],
		},
		{
			"ID #3",
			fields{
				testIds[2],
				0,
			},
			testIds[2],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &deltaUpdate{
				id:       tt.fields.id,
				progress: tt.fields.progress,
			}
			if got := update.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deltaUpdate.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deltaUpdate_NewVal(t *testing.T) {
	testIds := generateTestIds(6)
	type fields struct {
		id       uuid.UUID
		progress float32
	}
	type args struct {
		prev float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		{
			"To 100",
			fields{
				testIds[0],
				0.5,
			},
			args{
				0.5,
			},
			1,
		},
		{
			"To 100 overflow",
			fields{
				testIds[1],
				0.5,
			},
			args{
				1,
			},
			1,
		},
		{
			"To 50 from 0",
			fields{
				testIds[2],
				0.5,
			},
			args{
				0,
			},
			0.5,
		},
		{
			"To 50 from 25",
			fields{
				testIds[3],
				0.25,
			},
			args{0.25},
			0.5,
		},
		{
			"To 0",
			fields{
				testIds[4],
				-0.5,
			},
			args{
				0.5,
			},
			0,
		},
		{
			"To 0 overflow",
			fields{
				testIds[5],
				-1,
			},
			args{0.5},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &deltaUpdate{
				id:       tt.fields.id,
				progress: tt.fields.progress,
			}
			if got := update.NewVal(tt.args.prev); got != tt.want {
				t.Errorf("deltaUpdate.NewVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDeltaUpdate(t *testing.T) {
	testIds := generateTestIds(3)
	type args struct {
		id       uuid.UUID
		progress float32
	}
	tests := []struct {
		name string
		args args
		want *deltaUpdate
	}{
		{
			"No special case",
			args{
				testIds[0],
				0.5,
			},
			&deltaUpdate{
				testIds[0],
				0.5,
			},
		},
		{
			"Overflow",
			args{
				testIds[1],
				2,
			},
			&deltaUpdate{
				testIds[1],
				1,
			},
		},
		{
			"Underflow",
			args{
				testIds[2],
				-0.15,
			},
			&deltaUpdate{
				testIds[2],
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDeltaUpdate(tt.args.id, tt.args.progress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDeltaUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newFixedUpdate(t *testing.T) {
	testIds := generateTestIds(3)
	type args struct {
		id       uuid.UUID
		progress float32
	}
	tests := []struct {
		name string
		args args
		want *fixedUpdate
	}{
		{
			"No special case",
			args{
				testIds[0],
				0.5,
			},
			&fixedUpdate{
				testIds[0],
				0.5,
			},
		},
		{
			"Overflow",
			args{
				testIds[1],
				2,
			},
			&fixedUpdate{
				testIds[1],
				1,
			},
		},
		{
			"Underflow",
			args{
				testIds[2],
				-0.15,
			},
			&fixedUpdate{
				testIds[2],
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFixedUpdate(tt.args.id, tt.args.progress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFixedUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
