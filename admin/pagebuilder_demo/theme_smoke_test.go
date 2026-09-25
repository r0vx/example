package pagebuilder_demo

import (
	"os"
	"testing"

	"github.com/r0vx/admin/pagebuilder"
)

// TestDemoThemeIsComplete 示例主题必须提供它声明的全部模板，
// 否则示例本身就是坏的，主题作者照着抄会踩坑
func TestDemoThemeIsComplete(t *testing.T) {
	b := pagebuilder.New("/pages", nil).
		ThemeDev(true).
		RegisterTheme("demo", os.DirFS("../../themes/demo")).
		DefaultTheme("demo")

	registerDemoContainers(b)

	if err := b.ValidateThemes(); err != nil {
		t.Fatalf("示例主题不完整: %v", err)
	}
}

// TestDemoThemeFSWorksOutsideExampleRoot 工作目录不是 example 根（集成测试、
// 编译后的二进制放到别处）时读不到磁盘上的 themes/demo，必须回退到编译进二进制的副本，
// 否则 RegisterTheme 在配置期 panic，整个应用起不来
func TestDemoThemeFSWorksOutsideExampleRoot(t *testing.T) {
	b := pagebuilder.New("/pages", nil).RegisterTheme("demo", demoThemeFS())
	registerDemoContainers(b)
	if err := b.ValidateThemes(); err != nil {
		t.Fatalf("回退副本不可用: %v", err)
	}
}
