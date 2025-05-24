package pinyin

import (
	"testing"
	"reflect"
)

func TestNewArgs(t *testing.T) {
	args := NewArgs()
	if args.Style != Normal {
		t.Errorf("Expected default style Normal, got %d", args.Style)
	}
	if args.Heteronym != false {
		t.Error("Expected default heteronym false")
	}
}

func TestPinyin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		args     *Args
		expected [][]string
	}{
		{
			name:  "Normal style",
			input: "你好",
			args:  &Args{Style: Normal},
			expected: [][]string{{"ni"}, {"hao"}},
		},
		{
			name:  "Tone style",
			input: "你好",
			args:  &Args{Style: Tone},
			expected: [][]string{{"nǐ"}, {"hǎo"}},
		},
		{
			name:  "Heteronym",
			input: "了",
			args:  &Args{Heteronym: true},
			expected: [][]string{{"le", "liao", "liao"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pinyin(tt.input, tt.args)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		args     *Args
		expected string
	}{
		{
			name:     "Normal style",
			input:    "你好",
			args:     &Args{Style: Normal},
			expected: "nihao",
		},
		{
			name:     "First letter",
			input:    "你好",
			args:     &Args{Style: FirstLetter},
			expected: "nh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input, tt.args)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHomophones(t *testing.T) {
	input := "你好"
	args := &Args{Style: Normal}
	result := Homophones(input, args)

	expectedKeys := []string{"ni", "hao"}
	for _, key := range expectedKeys {
		if _, exists := result[key]; !exists {
			t.Errorf("Expected key %s not found", key)
		}
	}
}

func TestInitial(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"你好", "nh"},
		{"世界", "sj"},
		{"a", "a"}, // non-Chinese should pass through
	}

	for _, tt := range tests {
		result := Initial(tt.input)
		if result != tt.expected {
			t.Errorf("For %s, expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestNameVariations(t *testing.T) {
	args := &Args{
		Style:     Normal,
		Heteronym: true,
	}

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "单字多音",
			input: "了",
			expected: [][]string{{"le", "liao"}},
		},
		{
			name:  "多字组合",
			input: "重庆",
			expected: [][]string{{"zhong", "chong", "tong"}, {"qing"}},
		},
		{
			name:  "混合字符",
			input: "a你",
			expected: [][]string{{"a"}, {"ni"}},
		},
		{
			name:  "非中文字符",
			input: "abc",
			expected: [][]string{{"a"}, {"b"}, {"c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NameVariations(tt.input, args)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For %s, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
