package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			fields:  fields{},
			args:    args{filename: "./test/Parse_1.log"},
			wantErr: false,
		},
		{
			name:    "File not found",
			fields:  fields{},
			args:    args{filename: "./test/Parse_2.log"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			if err := p.Parse(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_parseLine(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			p.parseLine(tt.args.line)
		})
	}
}

func TestParser_gameKey(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Success",
			fields: fields{gameCounter: 5},
			want:   "game_05",
		},
		{
			name:   "Success",
			fields: fields{},
			want:   "game_00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			if got := p.gameKey(); got != tt.want {
				t.Errorf("gameKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_checkErrorState(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{
			name: "Error State = true",
			fields: fields{
				errorState:  true,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {},
				},
			},
			want: fields{
				errorState:  true,
				gameCounter: 0,
				Log:         map[string]Game{},
			},
		},
		{
			name: "Error State = false",
			fields: fields{
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {},
				},
			},
			want: fields{
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			p.checkErrorState()
			assert.Equal(t, p.Log, tt.want.Log)
		})
	}
}

func TestParser_initGame(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			if got := p.initGame(); got != tt.want {
				t.Errorf("initGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_addPlayer(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			if got := p.addPlayer(); got != tt.want {
				t.Errorf("addPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_addKill(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			if got := p.addKill(); got != tt.want {
				t.Errorf("addKill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_addWeaponKill(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	type args struct {
		weapon string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			p.addWeaponKill(tt.args.weapon)
		})
	}
}

func TestParser_addPlayerKill(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	type args struct {
		killer string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			p.addPlayerKill(tt.args.killer)
		})
	}
}

func TestParser_addWorldKill(t *testing.T) {
	type fields struct {
		line        string
		errorState  bool
		gameCounter int
		Log         map[string]Game
	}
	type args struct {
		victim string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				line:        tt.fields.line,
				errorState:  tt.fields.errorState,
				gameCounter: tt.fields.gameCounter,
				Log:         tt.fields.Log,
			}
			p.addWorldKill(tt.args.victim)
		})
	}
}
