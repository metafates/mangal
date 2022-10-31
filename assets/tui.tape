Require "mangal"
Output assets/tui.gif

Set FontSize 21
Set Width 1200
Set Height 1000

Hide
Type "mangal inline anilist set --name "Chainsaw Man" --id 138829" Enter
Type "rm -rf $(mangal where --downloads)/Chainsaw_Man" Enter
Type "clear" Enter
Show

Type@200ms "mangal" Sleep 200 ms Enter
Sleep 3s

# move cursor up and down
# select the first source
Down@400ms 2 Up@400ms 2 Enter
Sleep 2s

# start searching for "chain"
Type@200ms "chain" 
Sleep 1s 

# Accept auto-complete
Tab Sleep 1s Enter Sleep 4s

# Open filter
Type "/" Sleep 1s

# Filter out mangas that has "man" in their title
# Choose the first one
Type@200ms "man" Enter
Sleep 1s

# Open chapters
Enter

Sleep 5s

# Open anilist mangas that corresponds to the manga
Type "a" 
Sleep 4s
# Move cursor up
Type "k"
Sleep 2s
# Select it to bind with
Enter
Sleep 3s

# Mark chapter for download
Type " "
Sleep 1s
# Move cursor down and mark another chapter
Type "j"
Sleep 1s
Type " "
Sleep 1s

# Download
Sleep 2s
Enter

# Confirm download
Sleep 3s
Enter

# Wait for download to finish
Sleep 20s
