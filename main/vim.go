package main

import k "github.com/zyedidia/kbd"

func vim() k.Pattern {
	digit := k.RangeRune('0', '9')

	num := k.Cap(
		k.Plus(digit),
		"$0",
	)
	numopt := k.Cap(
		k.Star(digit),
		"$0",
	)

	move := k.Alt(
		k.Cap(k.Seq(numopt, k.MustLit("w")), "word-front -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("b")), "word-back -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("e")), "word-end -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("h")), "cursor-left -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("j")), "cursor-down -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("k")), "cursor-up -n $1"),
		k.Cap(k.Seq(numopt, k.MustLit("l")), "cursor-right -n $1"),
		k.Cap(k.MustLit("G"), "cursor-end-buffer"),
		k.Cap(k.Seq(num, k.MustLit("G")), "cursor-line-to $1"),
	)

	// repmove := k.Cap(
	// 	k.Seq(numopt, move),
	// 	"for { set i 0; set p $pos } { $i < $1 } { incr i } { set p [$2 $p] }",
	// )

	action := k.Alt(
		k.Cap(k.Seq(k.MustLit("Z"), k.MustLit("Z")), "save; quit"),
		k.Cap(move, "cursor-to [+ $pos [$1]]"),
	)

	raction := k.Alt(
		k.Cap(k.Seq(k.MustLit("d"), k.MustLit("d")), "delete-line"),
		k.Cap(k.Seq(k.MustLit("d"), move), "delete-range $pos [+ $pos [$1]]"),
		k.Cap(k.MustLit("D"), "exec 'd$'"),
	)

	bindings := k.Alt(
		action,
		k.Cap(k.Seq(numopt, raction), "repeat -n $1 $2"),
	)

	return k.Seq(bindings, k.End())
}
