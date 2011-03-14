function fixtype32() {
	if ($3 ~ /^-/)
		$3 = "(1<<32)" $3
}

/^const \(/ {
	cnst++
}
cnst && /^\)/ {
	cnst--
}


cnst && /MAXDWORD/ {
	fixtype32()
}
cnst && /HKEY_/ {
	fixtype32()
}

{
	print
}
