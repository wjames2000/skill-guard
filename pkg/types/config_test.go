package types

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Paths) != 1 || cfg.Paths[0] != "." {
		t.Errorf("默认路径应为 [.], 得到 %v", cfg.Paths)
	}
	if cfg.MaxSize != 10*1024*1024 {
		t.Errorf("默认 MaxSize 应为 10MB, 得到 %d", cfg.MaxSize)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"有效默认配置", DefaultConfig(), false},
		{"空路径", &Config{Paths: []string{}, MaxSize: 1}, true},
		{"无效严重级别", &Config{Paths: []string{"."}, Severity: "invalid", MaxSize: 1}, true},
		{"有效严重级别", &Config{Paths: []string{"."}, Severity: "high", MaxSize: 1}, false},
		{"MaxSize 为 0", &Config{Paths: []string{"."}, MaxSize: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
