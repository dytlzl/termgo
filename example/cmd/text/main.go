package main

import (
	"strings"

	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	quitMenu := new(QuitMenu)
	isOpenQuitMenu := false
	handleEvent := func(event any) any {
		switch typed := event.(type) {
		case rune:
			if isOpenQuitMenu {
				switch typed {
				case key.Esc:
					isOpenQuitMenu = false
				case key.ArrowLeft:
					quitMenu.IsYes = true
				case key.ArrowRight:
					quitMenu.IsYes = false
				case key.Enter:
					return tui.Terminate
				}
			} else {
				switch typed {
				case 'q', key.Esc:
					quitMenu.IsYes = false
					isOpenQuitMenu = true
				}
			}
		}
		return nil
	}
	err := tui.Run(func() *tui.View {
		return tui.ZStack(
			tui.InlineStack(
				tui.String(dograMagra1).Bold(),
				tui.String(dograMagra2).Italic(),
				tui.String(dograMagra3).Underline(),
				tui.String(dograMagra4).Strikethrough(),
			).RelativeSize(10, 10).Title("ドグラマグラ - 夢野久作").Border(),
			quitMenu.View().Invert(true).Hidden(!isOpenQuitMenu),
		).BGColor(color.RGB(145, 0, 145))
	},
		tui.OptionEventHandler(handleEvent),
	)
	if err != nil {
		panic(err)
	}
}

type QuitMenu struct {
	IsYes bool
}

func (q *QuitMenu) View() *tui.View {
	return tui.InlineStack(
		tui.Fmt("%sAre you sure to quit?\n\n%s     ",
			strings.Repeat(" ", (32-21)/2),
			strings.Repeat(" ", (32-21)/2)),
		tui.String(" Yes ").Invert(q.IsYes),
		tui.String(" "),
		tui.String(" No ").Invert(q.IsYes),
	).AbsoluteSize(36, 7).Title("Quit").Border()
}

const (
	dograMagra1 = ` …………ブウウ――――――ンンン――――――ンンンン………………。
 私がウスウスと眼を覚ました時、こうした蜜蜂（みつばち）の唸（うな）るような音は、まだ、その弾力の深い余韻を、私の耳の穴の中にハッキリと引き残していた。
 それをジッと聞いているうちに……今は真夜中だな……と直覚した。そうしてどこか近くでボンボン時計が鳴っているんだな……と思い思い、又もウトウトしているうちに、その蜜蜂のうなりのような余韻は、いつとなく次々に消え薄れて行って、そこいら中がヒッソリと静まり返ってしまった。
 私はフッと眼を開いた。
 かなり高い、白ペンキ塗の天井裏から、薄白い塵埃（ほこり）に蔽（おお）われた裸の電球がタッタ一つブラ下がっている。その赤黄色く光る硝子球（ガラスだま）の横腹に、大きな蠅（はえ）が一匹とまっていて、死んだように凝然（じっ）としている。その真下の固い、冷めたい人造石の床の上に、私は大の字型（なり）に長くなって寝ているようである。
`
	dograMagra2 = ` ……おかしいな…………。
 私は大の字型（なり）に凝然（じっ）としたまま、瞼（まぶた）を一パイに見開いた。そうして眼の球（たま）だけをグルリグルリと上下左右に廻転さしてみた。
 青黒い混凝土（コンクリート）の壁で囲まれた二間（けん）四方ばかりの部屋である。
 その三方の壁に、黒い鉄格子と、鉄網（かなあみ）で二重に張り詰めた、大きな縦長い磨硝子（すりガラス）の窓が一つ宛（ずつ）、都合三つ取付けられている、トテも要心（ようじん）堅固に構えた部屋の感じである。
 窓の無い側の壁の附け根には、やはり岩乗（がんじょう）な鉄の寝台が一個、入口の方向を枕にして横たえてあるが、その上の真白な寝具が、キチンと敷き展（なら）べたままになっているところを見ると、まだ誰も寝たことがないらしい。
`
	dograMagra3 = `	……おかしいぞ…………。
 私は少し頭を持ち上げて、自分の身体（からだ）を見廻わしてみた。
 白い、新しいゴワゴワした木綿の着物が二枚重ねて着せてあって、短かいガーゼの帯が一本、胸高に結んである。そこから丸々と肥（ふと）って突き出ている四本の手足は、全体にドス黒く、垢だらけになっている……そのキタナラシサ……。
`
	dograMagra4 = ` ……いよいよおかしい……。
 怖（こ）わ怖（ご）わ右手（めて）をあげて、自分の顔を撫（な）でまわしてみた。
 ……鼻が尖（と）んがって……眼が落ち窪（くぼ）んで……頭髪（あたま）が蓬々（ぼうぼう）と乱れて……顎鬚（あごひげ）がモジャモジャと延びて……。
 ……私はガバと跳ね起きた。
 モウ一度、顔を撫でまわしてみた。
 そこいらをキョロキョロと見廻わした。
 ……誰だろう……俺はコンナ人間を知らない……。
 胸の動悸がみるみる高まった。早鐘を撞（つ）くように乱れ撃ち初めた……呼吸が、それに連れて荒くなった。やがて死ぬかと思うほど喘（あえ）ぎ出した。……かと思うと又、ヒッソリと静まって来た。
 ……こんな不思議なことがあろうか……。
 ……自分で自分を忘れてしまっている……。
 ……いくら考えても、どこの何者だか思い出せない。……自分の過去の思い出としては、たった今聞いたブウ――ンンンというボンボン時計の音がタッタ一つ、記憶に残っている。……ソレッ切りである……。
 ……それでいて気は慥（たし）かである。森閑（しんかん）とした暗黒が、部屋の外を取巻いて、どこまでもどこまでも続き広がっていることがハッキリと感じられる……。
 ……夢ではない……たしかに夢では…………。
 私は飛び上った。
 ……窓の前に駈け寄って、磨硝子の平面を覗いた。そこに映った自分の容貌（かおかたち）を見て、何かの記憶を喚（よ）び起そうとした。……しかし、それは何にもならなかった。磨硝子の表面には、髪の毛のモジャモジャした悪鬼のような、私自身の影法師しか映らなかった。
 私は身を飜（ひるがえ）して寝台の枕元に在る入口の扉（ドア）に駈け寄った。鍵穴だけがポツンと開いている真鍮（しんちゅう）の金具に顔を近付けた。けれどもその金具の表面は、私の顔を写さなかった。只、黄色い薄暗い光りを反射するばかりであった。
 ……寝台の脚を探しまわった。寝具を引っくり返してみた。着ている着物までも帯を解いて裏返して見たけれども、私の名前は愚（おろ）か、頭文字らしいものすら発見し得なかった。
 私は呆然となった。私は依然として未知の世界に居る未知の私であった。私自身にも誰だかわからない私であった。
 こう考えているうちに、私は、帯を引きずったまま、無限の空間を、ス――ッと垂直に、どこへか落ちて行くような気がしはじめた。臓腑（はらわた）の底から湧き出して来る戦慄（せんりつ）と共に、我を忘れて大声をあげた。
 それは金属性を帯びた、突拍子（とっぴょうし）もない甲高（かんだか）い声であった……が……その声は私に、過去の何事かを思い出させる間もないうちに、四方のコンクリート壁に吸い込まれて、消え失せてしまった。
 又叫んだ。……けれども矢張（やは）り無駄であった。その声が一しきり烈（はげ）しく波動して、渦巻いて、消え去ったあとには、四つの壁と、三つの窓と、一つの扉が、いよいよ厳粛に静まり返っているばかりである。
 又叫ぼうとした。……けれどもその声は、まだ声にならないうちに、咽喉（のど）の奥の方へ引返してしまった。叫ぶたんびに深まって行く静寂の恐ろしさ……。
 奥歯がガチガチと音を立てはじめた。膝頭（ひざがしら）が自然とガクガクし出した。それでも自分自身が何者であったかを思い出し得ない……その息苦しさ。
 私は、いつの間にか喘（あえ）ぎ初めていた。叫ぼうにも叫ばれず、出ようにも出られぬ恐怖に包まれて、部屋の中央（まんなか）に棒立ちになったまま喘いでいた。
 ……ここは監獄か……精神病院か……。
 そう思えば思うほど高まる呼吸の音が、凩（こがらし）のように深夜の四壁に反響するのを聞いていた。
 そのうちに私は気が遠くなって来た。眼の前がズウ――と真暗くなって来た。そうして棒のように強直（ごうちょく）した全身に、生汗をビッショリと流したまま仰向（あおむ）け様（ざま）にスト――ンと、倒れそうになったので、吾知らず観念の眼を閉じた……と思ったが……又、ハッと機械のように足を踏み直した。両眼をカッと見開いて、寝台の向側の混凝土（コンクリート）壁を凝視した。
 その混凝土壁の向側から、奇妙な声が聞えて来たからであった。
 ……それは確かに若い女の声と思われた。けれども、その音調はトテも人間の肉声とは思えないほど嗄（しゃが）れてしまって、ただ、底悲しい、痛々しい響（ひびき）ばかりが、混凝土の壁を透して来るのであった。
「……お兄さま。お兄さま。お兄さまお兄さまお兄さまお兄さまお兄さま。……モウ一度……今のお声を……聞かしてエ――ッ…………」
 私は愕然（がくぜん）として縮み上った。思わずモウ一度、背後（うしろ）を振り返った。この部屋の中に、私以外の人間が一人も居ない事を承知し抜いていながら……それから又も、その女の声を滲（し）み透して来る、コンクリート壁の一部分を、穴のあく程、凝視した。
「……お兄さまお兄さまお兄さまお兄さまお兄さま……お隣りのお部屋に居らっしゃるお兄様……あたしです。妾（あたし）です。お兄様の許嫁（いいなずけ）だった……貴方（あなた）の未来の妻でした妾……あたしです。あたしです。どうぞ……どうぞ今のお声をモウ一度聞かして……聞かして頂戴……聞かして……聞かしてエ――ッ……お兄様お兄様お兄様お兄様……おにいさまア――ッ……」
 私は眼瞼（まぶた）が痛くなるほど両眼を見開いた。唇をアングリと開いた。その声に吸い付けられるようにヒョロヒョロと二三歩前に出た。そうして両手で下腹をシッカリと押え付けた。そのまま一心に混凝土（コンクリート）の壁を白眼（にら）み付けた。
 それは聞いている者の心臓を虚空に吊るし上げる程のモノスゴイ純情の叫びであった。臓腑をドン底まで凍らせずには措（お）かないくらいタマラナイ絶体絶命の声であった。……いつから私を呼び初めたかわからぬ……そうしてこれから先、何千年、何万年、呼び続けるかわからない真剣な、深い怨（うら）みの声であった。それが深夜の混凝土壁の向うから私？ を呼びかけているのであった。`
)
