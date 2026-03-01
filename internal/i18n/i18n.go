package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

// Translator handles internationalization
type Translator struct {
	translations map[string]string
	lang         string
	mu           sync.RWMutex
}

var (
	defaultTranslator *Translator
	once              sync.Once
)

// Language codes
const (
	LangEnglish = "en"
	LangChinese = "zh"
)

// Init initializes the i18n system with the given language
func Init(lang string) {
	once.Do(func() {
		defaultTranslator = &Translator{
			lang: lang,
		}
		defaultTranslator.loadTranslations(lang)
	})
}

// GetLanguage returns the current language
func GetLanguage() string {
	if defaultTranslator == nil {
		return LangEnglish
	}
	return defaultTranslator.lang
}

// SetLanguage changes the current language
func SetLanguage(lang string) {
	if defaultTranslator != nil {
		defaultTranslator.mu.Lock()
		defaultTranslator.lang = lang
		defaultTranslator.mu.Unlock()
		defaultTranslator.loadTranslations(lang)
	}
}

// DetectSystemLanguage detects the system language from environment
func DetectSystemLanguage() string {
	// Check LANG environment variable (Linux/macOS)
	lang := os.Getenv("LANG")
	if lang != "" {
		// Parse language code (e.g., "zh_CN.UTF-8" -> "zh")
		if strings.HasPrefix(lang, "zh") {
			return LangChinese
		}
		if strings.HasPrefix(lang, "en") {
			return LangEnglish
		}
	}

	// Check LC_ALL
	lcAll := os.Getenv("LC_ALL")
	if lcAll != "" {
		if strings.HasPrefix(lcAll, "zh") {
			return LangChinese
		}
		if strings.HasPrefix(lcAll, "en") {
			return LangEnglish
		}
	}

	// Check LANGUAGE
	language := os.Getenv("LANGUAGE")
	if language != "" {
		if strings.Contains(language, "zh") {
			return LangChinese
		}
	}

	// Check macOS system preferences (if LANG not set)
	runtimeOS := runtime.GOOS
	if runtimeOS == "darwin" {
		// Try to read macOS system locale
		if output, err := exec.Command("defaults", "read", "-g", "AppleLocale").Output(); err == nil {
			locale := strings.TrimSpace(string(output))
			if strings.HasPrefix(locale, "zh") {
				return LangChinese
			}
			if strings.HasPrefix(locale, "en") {
				return LangEnglish
			}
		}
		// Try AppleLanguages as fallback
		if output, err := exec.Command("defaults", "read", "-g", "AppleLanguages").Output(); err == nil {
			if strings.Contains(string(output), "zh") {
				return LangChinese
			}
		}
	}

	// Default to English
	return LangEnglish
}

// T translates a key with optional arguments
func T(key string, args ...interface{}) string {
	if defaultTranslator == nil {
		Init(DetectSystemLanguage())
	}
	return defaultTranslator.Translate(key, args...)
}

// Translate translates a key with optional arguments
func (t *Translator) Translate(key string, args ...interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.translations == nil {
		return key
	}

	template, ok := t.translations[key]
	if !ok {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(template, args...)
	}

	return template
}

// loadTranslations loads translations from embedded files
func (t *Translator) loadTranslations(lang string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	filename := fmt.Sprintf("locales/%s.json", lang)
	data, err := localesFS.ReadFile(filename)
	if err != nil {
		// Fallback to English if requested language not found
		if lang != LangEnglish {
			data, err = localesFS.ReadFile("locales/en.json")
			if err != nil {
				t.translations = make(map[string]string)
				return
			}
		} else {
			t.translations = make(map[string]string)
			return
		}
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		t.translations = make(map[string]string)
		return
	}

	t.translations = translations
}
