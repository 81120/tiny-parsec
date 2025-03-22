package parser_test

import (
	"strings"
	"testing"

	. "github.com/81120/tiny-parsec/parser"
	"github.com/stretchr/testify/assert"
)

func TestFmap(t *testing.T) {
	t.Run("正常转换", func(t *testing.T) {
		p := Fmap(Str("hello"), strings.ToUpper)
		result := p.Parse("hello world")
		assert.True(t, result.IsJust())
		assert.Equal(t, "HELLO", result.Get().First)
		assert.Equal(t, " world", result.Get().Second)
	})

	t.Run("空输入处理", func(t *testing.T) {
		p := Fmap(Str("hello"), strings.ToUpper)
		result := p.Parse("")
		assert.True(t, result.IsNothing())
	})
}

func TestBind(t *testing.T) {
	t.Run("序列解析成功", func(t *testing.T) {
		p := Bind(Str("user:"), func(_ string) Parser[string] {
			return Fmap(Alphas(), strings.ToLower)
		})
		result := p.Parse("user:ADMIN")
		assert.True(t, result.IsJust())
		assert.Equal(t, "admin", result.Get().First)
	})

	t.Run("中间失败处理", func(t *testing.T) {
		p := Bind(Str("123"), func(_ string) Parser[string] { return Str("abc") })
		result := p.Parse("124abc")
		assert.True(t, result.IsNothing())
	})
}

func TestOrElse(t *testing.T) {
	t.Run("优先匹配成功", func(t *testing.T) {
		p := OrElse(Str("a"), Str("b"))
		result := p.Parse("a")
		assert.True(t, result.IsJust())
		assert.Equal(t, "a", result.Get().First)
	})

	t.Run("备选匹配成功", func(t *testing.T) {
		p := OrElse(Str("a"), Str("b"))
		result := p.Parse("b")
		assert.True(t, result.IsJust())
		assert.Equal(t, "b", result.Get().First)
	})
}

func TestZeroOrOne(t *testing.T) {
	t.Run("零次匹配", func(t *testing.T) {
		p := ZeroOrOne(Str("optional"))
		result := p.Parse("")
		assert.True(t, result.IsJust())
		assert.True(t, result.Get().First.IsNothing())
	})

	t.Run("一次匹配", func(t *testing.T) {
		p := ZeroOrOne(Str("item"))
		result := p.Parse("item123")
		assert.True(t, result.IsJust())
		assert.Equal(t, "item", result.Get().First.Get())
		assert.Equal(t, "123", result.Get().Second)
	})
}

func TestZeroOrMore(t *testing.T) {
	t.Run("零次匹配", func(t *testing.T) {
		p := ZeroOrMore(Str("a"))
		result := p.Parse("b")
		assert.True(t, result.IsJust())
		assert.Empty(t, result.Get().First)
	})

	t.Run("多次匹配", func(t *testing.T) {
		p := ZeroOrMore(Str("a"))
		result := p.Parse("aaaaab")
		assert.Len(t, result.Get().First, 5)
		assert.Equal(t, "b", result.Get().Second)
	})
}

func TestOneOrMore(t *testing.T) {
	t.Run("至少一次匹配", func(t *testing.T) {
		p := OneOrMore(Str("a"))
		result := p.Parse("aab")
		assert.True(t, result.IsJust())
		assert.Len(t, result.Get().First, 2)
	})

	t.Run("零次失败", func(t *testing.T) {
		p := OneOrMore(Str("a"))
		result := p.Parse("b")
		assert.True(t, result.IsNothing())
	})
}

func TestTrim(t *testing.T) {
	t.Run("完整修剪效果", func(t *testing.T) {
		p := Trim(Str("core"))
		result := p.Parse("  \tcore\n  remaining")
		assert.True(t, result.IsJust())
		assert.Equal(t, "core", result.Get().First)
		assert.Equal(t, "remaining", result.Get().Second)
	})

	t.Run("空输入处理", func(t *testing.T) {
		p := Trim(Str("test"))
		result := p.Parse("  \t\n  ")
		assert.True(t, result.IsNothing())
	})
}

func TestTrimCombinators(t *testing.T) {
	t.Run("TrimLeft去除头部空白", func(t *testing.T) {
		p := TrimLeft(Str("data"))
		result := p.Parse("   data")
		assert.True(t, result.IsJust())
		assert.Equal(t, "data", result.Get().First)
	})

	t.Run("TrimRight保留尾部内容", func(t *testing.T) {
		p := TrimRight(Str("value"))
		result := p.Parse("value  	remaining")
		assert.True(t, result.IsJust())
		assert.Equal(t, "value", result.Get().First)
		assert.Equal(t, "remaining", result.Get().Second)
	})
}

func TestSepBy(t *testing.T) {
	t.Run("带分隔符的列表解析", func(t *testing.T) {
		p := SepBy(Str("item"), Char(','))
		result := p.Parse("item,item,item")
		assert.True(t, result.IsJust())
		assert.Len(t, result.Get().First, 3)
	})

	t.Run("空列表处理", func(t *testing.T) {
		p := SepBy(Str("item"), Char(','))
		result := p.Parse("")
		assert.True(t, result.IsJust())
		assert.Empty(t, result.Get().First)
	})
}

func TestSatisfyCombinators(t *testing.T) {
	t.Run("Satisfy条件匹配", func(t *testing.T) {
		p := Satisfy(func(r rune) bool { return r == 'X' })
		result := p.Parse("X123")
		assert.True(t, result.IsJust())
		assert.Equal(t, 'X', result.Get().First)
	})

	t.Run("SatisfyWith过滤条件", func(t *testing.T) {
		p := SatisfyWith(Str("123"), func(s string) bool { return s == "123" })
		result := p.Parse("1234")
		assert.True(t, result.IsJust())
		assert.Equal(t, "123", result.Get().First)
	})
}

func TestBetween(t *testing.T) {
	t.Run("包围结构解析", func(t *testing.T) {
		p := Between(Str("("), Str("content"), Str(")"))
		result := p.Parse("(content)rest")
		assert.True(t, result.IsJust())
		assert.Equal(t, "content", result.Get().First)
	})
}

func TestSeq(t *testing.T) {
	t.Run("顺序解析成功", func(t *testing.T) {
		p := Seq(Str("a"), Str("b"), Str("c"))
		result := p.Parse("abc")
		assert.True(t, result.IsJust())
		assert.Equal(t, []string{"a", "b", "c"}, result.Get().First)
	})

	t.Run("中间失败处理", func(t *testing.T) {
		p := Seq(Str("a"), Str("x"), Str("c"))
		result := p.Parse("abc")
		assert.True(t, result.IsNothing())
	})
}

func TestLazy(t *testing.T) {
	t.Run("延迟解析验证", func(t *testing.T) {
		called := false
		p := Lazy(func() Parser[string] {
			called = true
			return Str("lazy")
		})
		result := p.Parse("lazy")
		assert.True(t, called)
		assert.True(t, result.IsJust())
	})
}

func TestOmitSides(t *testing.T) {
	t.Run("OmitLeft保留右边", func(t *testing.T) {
		p := OmitLeft(Str("{"), Str("value"))
		result := p.Parse("{value}")
		assert.True(t, result.IsJust())
		assert.Equal(t, "value", result.Get().First)
		assert.Equal(t, "}", result.Get().Second)
	})

	t.Run("OmitRight保留左边", func(t *testing.T) {
		p := OmitRight(Str("key:"), Str(" "))
		result := p.Parse("key: value")
		assert.True(t, result.IsJust())
		assert.Equal(t, "key:", result.Get().First)
		assert.Equal(t, "value", result.Get().Second)
	})
}
