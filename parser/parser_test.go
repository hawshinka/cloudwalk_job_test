package parser

import (
	"io/fs"
	"syscall"
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
	tests := []struct {
		name     string
		filename string
		fields   fields
		wantErr  error
	}{
		{
			name:     "Success",
			filename: "./test/Parse_1.log",
			fields:   fields{},
			wantErr:  nil,
		},
		{
			name:     "File not found",
			filename: "./test/Parse_2.log",
			fields:   fields{},
			wantErr:  &fs.PathError{Op: "open", Path: "./test/Parse_2.log", Err: syscall.ENOENT},
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
			err := p.Parse(tt.filename)
			assert.Equal(t, tt.wantErr, err)
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
			key := p.gameKey()
			assert.Equal(t, tt.want, key)
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
				gameCounter: 1,
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
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
		want   fields
	}{
		{
			name: "Not InitGame line",
			fields: fields{
				line:        " 15:00 Exit: Timelimit hit.",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        " 15:00 Exit: Timelimit hit.",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Success",
			fields: fields{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Players:      make([]string, 0),
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
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
			p.initGame()
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
		want   fields
	}{
		{
			name:   "Not ClientUserinfoChanged line",
			fields: fields{},
			want:   fields{},
		},
		{
			name: "Error in regex",
			fields: fields{
				line:        " 20:38 ClientUserinfoChanged: 2",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        " 20:38 ClientUserinfoChanged: 2",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Game in error state",
			fields: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Player already exists in the list",
			fields: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
		},
		{
			name: "New player to the list",
			fields: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Dono da bola\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: fields{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Dono da bola\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido", "Dono da bola"},
					},
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
			p.addPlayer()
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
		want   fields
	}{
		{
			name:   "Not Kill line",
			fields: fields{},
			want:   fields{},
		},
		{
			name: "Error in regex",
			fields: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Game in error state",
			fields: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Success with player kill",
			fields: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: fields{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills: 1,
						Kills: map[string]int{
							"Isgalamido": 1,
						},
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
		},
		{
			name: "Success with world kill",
			fields: fields{
				line:        "  3:27 Kill: 1022 3 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: fields{
				line:        "  3:27 Kill: 1022 3 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills: 1,
						Kills: map[string]int{
							"Isgalamido": -1,
						},
						KillsByMeans: map[string]int{
							"MOD_TRIGGER_HURT": 1,
						},
					},
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
			p.addKill()
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
		weapon string
		fields fields
		want   fields
	}{
		{
			name:   "Success",
			weapon: "MOD_ROCKET",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
		},
		{
			name:   "Success with same weapon",
			weapon: "MOD_ROCKET",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 2,
						},
					},
				},
			},
		},
		{
			name:   "Success with different weapon",
			weapon: "MOD_ROCKET_SPLASH",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET":        1,
							"MOD_ROCKET_SPLASH": 1,
						},
					},
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
			p.addWeaponKill(tt.weapon)
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
	tests := []struct {
		name   string
		killer string
		fields fields
		want   fields
	}{
		{
			name:   "Success",
			killer: "Isgalamido",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a second hit from the same player",
			killer: "Isgalamido",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 2,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit from a different player",
			killer: "Zeh",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
							"Zeh":        1,
						},
					},
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
			p.addPlayerKill(tt.killer)
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
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
	tests := []struct {
		name   string
		victim string
		fields fields
		want   fields
	}{
		{
			name:   "Success",
			victim: "Isgalamido",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a second hit to the same player",
			victim: "Isgalamido",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -2,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit to a different player",
			victim: "Zeh",
			fields: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: fields{
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
							"Zeh":        -1,
						},
					},
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
			p.addWorldKill(tt.victim)
			assert.Equal(t, tt.want.line, p.line)
			assert.Equal(t, tt.want.errorState, p.errorState)
			assert.Equal(t, tt.want.gameCounter, p.gameCounter)
			assert.Equal(t, tt.want.Log, p.Log)
		})
	}
}
