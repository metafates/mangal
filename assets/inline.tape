Require "mangal"
Require "bat"
Require "jq"
Output assets/inline.gif


Set FontSize 21
Set Width 1200
Set Height 1000

Type@20ms "# see 'mangal inline schema' to view output schema" Enter
Type@70ms "mangal inline -S ComicK -q 'chainsaw' --manga first --json --include-anilist-manga | jq . | bat -l json"
Sleep 1s
Enter
Sleep 3s
Type@200ms "ddddddd"
Sleep 2
Type "q"

Enter

Type@20ms "# see 'mangal inline schema --anilist' to view output schema" Enter
Type@70ms "mangal inline anilist search --name 'chainsaw man' | jq . | bat -l json"
Sleep 1s
Enter
Sleep 3s
Type@200ms "ddddddd"
Sleep 2
