package gobas

import "strings"

/*
Quote from https://www.c64-wiki.de/wiki/PRINT

Folgende Eigenschaften und Besonderheiten weist PRINT auf:

* Wird Text ohne Anführungszeichen eingegeben, wird dieser als Ausdruck interpretiert, wobei häufig eine Zahlenvariable (falls mit einem Buchstaben beginnend) oder ein numerischer Wert (falls mit einer Ziffer beginnend) vorliegt.
* Der PRINT-Befehl schließt die Ausgaben immer mit einem Zeilenwechsel ab, es sei denn, der Ausdruck endet mit einem Semikolon (;) oder Komma (,). Der PRINT-Befehl ohne Parameterangabe bewirkt nur ein Zeilenwechsel.
* Steht hinter dem letzten Ausdrücken ein Semikolon (;), so erfolgt die nächste Bildschirmausgabe direkt hinter dem ausgegebenen Ausdruck. Für die Ausgabe von Zahlen (als Ergebnis eines numerischen Ausdrucks) gilt dabei:
	1. Bei positiven Werten erscheint ein führendes Leerzeichen (als Platzhalter für das Vorzeichen).
	2. Der Zahl folgt eine Leerstelle (bei Bildschirmausgabe genaugenommen ein {CRSR RIGHT}-Zeichen, sonst ein normales Leerzeichen).
* Mit einem Komma (,) kann jeweils die nächste Tabulatorposition (alle 10 Zeichen) angesprungen werden. Beispielsweise wird der Cursor auf Spalte 10 gesetzt, wenn er sich vorher in Spalte 0-9 befand. Es gelten die gleichen Aussagen bei den Zahlenausdrücken wie beim Semikolon.
* Das Setzen von Komma (,) oder Semikolon (;) ist nicht erforderlich, um Ausdrücke abzutrennen. Der BASIC-Interpreter verlangt diese im Zusammenhang mit PRINT nicht.
	Bei fehlendem Trennzeichen tritt das Verhalten eines Semikolons (;) ein.
	Jedoch werden Leerzeichen außerhalb von Anführungszeichen (" ") durch den BASIC-Interpreter (siehe CHRGET) ignoriert,
	weshalb aufeinander folgende Zahlenvariablen oder numerische Werte auf jeden Fall mit einem Trennzeichen getrennt werden müssen.
	Geschieht dies nicht, werden die Zeichen jeweils als einziger Variablenname bzw. numerischer Wert interpretiert.
* Die speziellen Funktionen SPC und TAB (Ausgabefunktionen) können nur hier (und bei den verwandten Kommandos PRINT# und CMD) verwendet werden. Sie haben am Ende der Ausgabe, abgesehen ihrer eigenen Funktion, das gleiche Verhalten wie ein Semikolon (;), d.h. sie unterdrücken den Zeilenwechsel.
* Mittels BASIC-Befehl CMD kann die Bildschirmausgabe mit PRINT auch auf ein anderes Gerät umgeleitet werden.
*/

//TODO: Print!

type printItem interface{}

type printSemicolon struct{}
type printComma struct{}

func mustParsePrint(raw string) PRINT {
	print := PRINT{}
	var curr string
	flush := func() {
		curr = strings.TrimSpace(curr)
		if curr == "" {
			return
		}
		expr := mustParseExpression(curr)
		print.Items = append(print.Items, expr)
		curr = ""
	}

	var inQuotes bool
	var bcount int
	for _, r := range raw {
		if r == '"' {
			if !inQuotes {
				flush()
				inQuotes = true
			} else {
				inQuotes = false
				curr += string(r)
				flush()
				continue
			}
		}
		if inQuotes {
			curr += string(r)
			continue
		}
		if r == '(' {
			bcount++
		}
		if r == ')' {
			bcount--
		}
		if bcount == 0 && r == ';' {
			flush()
			print.Items = append(print.Items, printSemicolon{})
			continue
		}
		if bcount == 0 && r == ',' {
			flush()
			print.Items = append(print.Items, printComma{})
			continue
		}
		curr += string(r)
	}
	flush()

	return print
}
