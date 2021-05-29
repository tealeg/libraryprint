package main

import (
	"fmt"
	"io"
	"strings"
)

func escapeLatex(s string) string {
	r := strings.ReplaceAll
	s = r(s, "\\", "\\textasciibackslash")
	s = r(s, "&", "\\&")
	s = r(s, "%", "\\%")
	s = r(s, "#", "\\#")
	s = r(s, "$", "\\$")
	s = r(s, "#", "\\#")
	s = r(s, "_", "\\_")
	s = r(s, "{", "\\{")
	s = r(s, "}", "\\}")
	s = r(s, "~", "\\textasciitilde")
	s = r(s, "^", "\\textasciicircum")
	return s
}

func writePreamble(w io.Writer) error {
	_, err := io.WriteString(w, `
\documentclass[10pt,a5paper]{article}
\usepackage[margin=0.5cm,footskip=0.5cm]{geometry}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{polyglossia}
\usepackage{xunicode}
\usepackage{amsmath}
\usepackage{amsfonts}
\usepackage{amssymb}
\usepackage{graphicx}
\usepackage{fontspec}
\pagenumbering{gobble}
\setmainfont{IBM Plex Serif}
\setmonofont{IBM Plex Mono}
\let\centering\relax

\begin{document}
`)
	return err
}


func writePostscript(w io.Writer) error {
	_, err := io.WriteString(w, `
\end{document}
`)
	return err
}


func writeCard(w io.Writer, lr *LibraryRecord) error {
	card := fmt.Sprintf(`
\begin{tabular}{lp{0.6\textwidth}}
  \texttt{\textbf{\Large TITEL}} & {\Large %s} \\
  \texttt{\textbf{AUTOR}}	& %s \\
  \texttt{\textbf{KATEGORIEN}} & %s \\
  \texttt{\textbf{ARCHIEVNUMMER}} & %s \\
\end{tabular}
\\
\\
\\
\begin{tabular}{p{0.45\textwidth}|p{0.45\textwidth}}
  \texttt{\textbf{AUSGELIEHEN VON}}	& \texttt{\textbf{AUSGELIEHEN AM}} \\
  \hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline
& \\& \\
\hline

\end{tabular}
\newpage
`,
		escapeLatex(lr.Titel),
		escapeLatex(lr.Autoren),
		escapeLatex(lr.Kategorien),
		escapeLatex(lr.ArchiveNummer),
	)
	_, err := io.WriteString(w, card)
	return err
}
