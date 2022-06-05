package example

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

var focus EventHandler = nil

func Run() {
	cc := new(ColorCode)
	kc := new(KeyCode)
	n := new(Note)
	q := new(Quit)
	hasFocus := func(handler EventHandler) tui.Style {
		if focus == handler {
			return focusedStyle
		} else {
			return style
		}
	}
	handleEvent := func(event any) any {
		if focus == nil {
			switch typed := event.(type) {
			case rune:
				switch typed {
				case 'c':
					focus = cc
				case 'k':
					focus = kc
				case 'n':
					focus = n
				case 'q', key.Esc:
					q.IsYes = false
					focus = q
				}
			}
		} else {
			switch typed := event.(type) {
			case rune:
				switch typed {
				case key.Esc:
					focus = nil
					return nil
				}
			}
			return focus.HandleEvent(event)
		}
		return nil
	}
	err := tui.Run(func() *tui.View {
		return tui.ZStack(
			tui.VStack(
				tui.HStack(
					tui.TextView(dograMagra).RelativeSize(7, 12).Title("Dogra Magra").Style(style).Border(style),
					tui.ViewWithRenderer(n.Body).RelativeSize(5, 12).Title("Note").Style(hasFocus(n)).Border(hasFocus(n)),
				).RelativeSize(12, 6),
				tui.HStack(
					tui.ViewWithRenderer(cc.Body).AbsoluteSize(69, 0).Title("Color Code").Style(hasFocus(cc)).Border(hasFocus(cc)),
					tui.ViewWithRenderer(kc.Body).Title("Key Code").Style(hasFocus(kc)).Border(hasFocus(kc)),
					tui.VStack(
						tui.HStack(
							tui.VStack(
								tui.HStack(
									tui.VStack().Border(style),
									tui.VStack().Border(style),
								),
								tui.HStack().Border(style),
							),
							tui.VStack().Border(style),
						),
						tui.HStack().Border(style),
					).Border(style).Title("Layout"),
				).Border(style),
				tui.TextView("Footer is here.").Style(style.Invert()).AbsoluteSize(0, 1).Padding(0, 1, 0, 0),
			).Title("Example"),
			tui.ViewWithRenderer(q.Body).AbsoluteSize(36, 7).Title("Quit").Style(style.Invert()).Border(style.Invert()).Hidden(focus != q),
		)
	},
		tui.OptionEventHandler(handleEvent),
		tui.OptionStyle(style),
	)
	if err != nil {
		panic(err)
	}
}

type EventHandler interface {
	HandleEvent(event any) any
}

var focusedStyle = tui.Style{F256: 255, B256: 53}
var style = tui.Style{F256: 218, B256: 53}

const dograMagra = "…………ブウウ――――――ンンン――――――ンンンン………………。\n私がウスウスと眼を覚ました時、こうした蜜蜂みつばちの唸うなるような音は、まだ、その弾力の深い余韻を、私の耳の穴の中にハッキリと引き残していた。\nそれをジッと聞いているうちに……今は真夜中だな……と直覚した。そうしてどこか近くでボンボン時計が鳴っているんだな……と思い思い、又もウトウトしているうちに、その蜜蜂のうなりのような余韻は、いつとなく次々に消え薄れて行って、そこいら中がヒッソリと静まり返ってしまった。\n私はフッと眼を開いた。\nかなり高い、白ペンキ塗の天井裏から、薄白い塵埃ほこりに蔽おおわれた裸の電球がタッタ一つブラ下がっている。その赤黄色く光る硝子球ガラスだまの横腹に、大きな蠅はえが一匹とまっていて、死んだように凝然じっとしている。その真下の固い、冷めたい人造石の床の上に、私は大の字型なりに長くなって寝ているようである。\n……おかしいな…………。\n私は大の字型なりに凝然じっとしたまま、瞼まぶたを一パイに見開いた。そうして眼の球たまだけをグルリグルリと上下左右に廻転さしてみた。\n青黒い混凝土コンクリートの壁で囲まれた二間けん四方ばかりの部屋である。\nその三方の壁に、黒い鉄格子と、鉄網かなあみで二重に張り詰めた、大きな縦長い磨硝子すりガラスの窓が一つ宛ずつ、都合三つ取付けられている、トテも要心ようじん堅固に構えた部屋の感じである。\n窓の無い側の壁の附け根には、やはり岩乗がんじょうな鉄の寝台が一個、入口の方向を枕にして横たえてあるが、その上の真白な寝具が、キチンと敷き展ならべたままになっているところを見ると、まだ誰も寝たことがないらしい。\n……おかしいぞ…………。\n私は少し頭を持ち上げて、自分の身体からだを見廻わしてみた。\n白い、新しいゴワゴワした木綿の着物が二枚重ねて着せてあって、短かいガーゼの帯が一本、胸高に結んである。そこから丸々と肥ふとって突き出ている四本の手足は、全体にドス黒く、垢だらけになっている……そのキタナラシサ……。\n……いよいよおかしい……。\n怖こわ怖ごわ右手めてをあげて、自分の顔を撫なでまわしてみた。\n……鼻が尖とんがって……眼が落ち窪くぼんで……頭髪あたまが蓬々ぼうぼうと乱れて……顎鬚あごひげがモジャモジャと延びて……。\n……私はガバと跳ね起きた。\nモウ一度、顔を撫でまわしてみた。\nそこいらをキョロキョロと見廻わした。\n……誰だろう……俺はコンナ人間を知らない……。\n胸の動悸がみるみる高まった。早鐘を撞つくように乱れ撃ち初めた……呼吸が、それに連れて荒くなった。やがて死ぬかと思うほど喘あえぎ出した。……かと思うと又、ヒッソリと静まって来た。\n……こんな不思議なことがあろうか……。\n……自分で自分を忘れてしまっている……。\n……いくら考えても、どこの何者だか思い出せない。……自分の過去の思い出としては、たった今聞いたブウ――ンンンというボンボン時計の音がタッタ一つ、記憶に残っている。……ソレッ切りである……。\n……それでいて気は慥たしかである。森閑しんかんとした暗黒が、部屋の外を取巻いて、どこまでもどこまでも続き広がっていることがハッキリと感じられる……。\n……夢ではない……たしかに夢では…………。\n私は飛び上った。\n……窓の前に駈け寄って、磨硝子の平面を覗いた。そこに映った自分の容貌かおかたちを見て、何かの記憶を喚よび起そうとした。……しかし、それは何にもならなかった。磨硝子の表面には、髪の毛のモジャモジャした悪鬼のような、私自身の影法師しか映らなかった。"

type Quit struct {
	IsYes bool
}

func (q *Quit) Body(s tui.Size) []tui.Text {
	style := style.Invert()
	t := []tui.Text{
		{Str: strings.Repeat("\n", (s.Height-3)/2), Style: style},
		{Str: strings.Repeat(" ", (s.Width-21)/2) + "Are you sure to quit?\n\n", Style: style},
		{Str: strings.Repeat(" ", (s.Width-21)/2) + "     ", Style: style},
		{Str: " Yes ", Style: style},
		{Str: "  ", Style: style},
		{Str: " No ", Style: style},
	}
	if q.IsYes {
		t[3].Style = style.Invert()
	} else {
		t[5].Style = style.Invert()
	}
	return t
}

func (q *Quit) HandleEvent(event any) any {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.Enter:
			if q.IsYes {
				return tui.Terminate
			} else {
				focus = nil
			}
		case key.ArrowLeft:
			q.IsYes = true
		case key.ArrowRight:
			q.IsYes = false
		}
	}
	return nil
}

type Note struct {
	position int
	input    string
}

func (n *Note) Body(tui.Size) []tui.Text {
	style := tui.Style{F256: 255, B256: 53}
	cursorStyle := tui.Style{F256: style.B256, B256: style.F256, HasCursor: true}
	if focus == n {
		if n.position == len(n.input) {
			return []tui.Text{
				{Str: n.input[:n.position], Style: style},
				{Str: " ", Style: cursorStyle},
			}
		}
		r, size := utf8.DecodeRuneInString(n.input[n.position:])
		if r == '\n' {
			return []tui.Text{
				{Str: n.input[:n.position], Style: style},
				{Str: " ", Style: cursorStyle},
				{Str: n.input[n.position:], Style: tui.Style{F256: 255, B256: 53}},
			}
		}
		return []tui.Text{
			{Str: n.input[:n.position], Style: style},
			{Str: n.input[n.position : n.position+size], Style: cursorStyle},
			{Str: n.input[n.position+size:], Style: style},
		}
	} else {
		return []tui.Text{
			{Str: n.input, Style: style},
		}
	}
}

func (n *Note) HandleEvent(event any) any {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.Enter:
			n.input += "\n"
			n.position++
		case key.ArrowLeft:
			if n.position > 0 {
				_, size := utf8.DecodeLastRuneInString(n.input[:n.position])
				n.position -= size
			}
		case key.ArrowRight:
			if n.position < len(n.input) {
				_, size := utf8.DecodeRuneInString(n.input[n.position:])
				n.position += size
			}
		case key.ArrowUp, key.ArrowDown:
		case key.Del:
			if n.input != "" {
				_, size := utf8.DecodeLastRuneInString(n.input[:n.position])
				n.input = n.input[:n.position-size] + n.input[n.position:]
				n.position -= size
			}
		default:
			n.input = n.input[:n.position] + string(typed) + n.input[n.position:]
			n.position += utf8.RuneLen(typed)
		}
	}
	return nil
}

type ColorCode struct {
	position int
}

func (m *ColorCode) HandleEvent(event any) any {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.ArrowLeft:
			if m.position > 0 {
				m.position--
			}
		case key.ArrowRight:
			if m.position < 255 {
				m.position++
			}
		case key.ArrowUp:
			if m.position > 15 {
				m.position -= 16
			}
		case key.ArrowDown:
			if m.position < 240 {
				m.position += 16
			}
		}
	}
	return nil
}

func (m *ColorCode) Body(tui.Size) []tui.Text {
	var slice []tui.Text
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			seq := i*16 + j
			if seq == m.position {
				slice = append(slice, tui.Text{Str: " ", Style: style})
				slice = append(slice, tui.Text{Str: fmt.Sprintf("%3d", i*16+j), Style: tui.Style{F256: seq, B256: 103}})
			} else {
				slice = append(slice, tui.Text{Str: fmt.Sprintf("%4d", i*16+j), Style: tui.Style{F256: seq, B256: style.B256}})
			}
		}
		slice = append(slice, tui.Text{Str: "\n", Style: style})
	}
	return slice
}

type KeyCode struct {
	codes []rune
}

func (k *KeyCode) Body(tui.Size) []tui.Text {
	slice := make([]tui.Text, 0, len(k.codes))
	for i := len(k.codes) - 1; i >= 0 && i > len(k.codes)-5000; i-- {
		slice = append(slice, tui.Text{Str: fmt.Sprintf(" %d", int(k.codes[i])), Style: style})
	}
	return slice
}

func (k *KeyCode) HandleEvent(event any) any {
	if k.codes == nil {
		k.codes = make([]rune, 0)
	}
	switch typed := event.(type) {
	case rune:
		k.codes = append(k.codes, typed)
	}
	return nil
}
