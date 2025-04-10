# go-reloaded

Purpose: The tool edits, completes, auto-corrects the text

Input: the first argument <sample.txt> file with a text to modify <br> Output: the second argument <result.txt> file with modified text

### Formatting text
Purpose: program formats text into uppercase, capitalized, lowcase, binary and hexadecimal versions <br> Input: text from <sample.txt> file <br> Output: formatted text printed to <result.txt> file  <br> 

- Block for text formatting

**makeCases** finds **(low|up|cap|bin|hex, number)** by regex `\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)` and recursively calls processText for further formatting <br>

**processText** formats text by one match, where ***leftText***, a substring extracted from text by `regex.FindAllStringIndex` and ***data[i][2]***, number to modify extracted by `regex.FindAllStringSubmatch`, if number is invalid it moves to the next iteration, removing invalid match, if valid match it modifies the text and returns a formatted text printing errors and invalid matches. Once makeCases receives a flag it stops recursively calling <br> `In case number is above length of words in leftText, it returns error that is invalid match` <br>`In case of invalid number, that match is considered invalid and collected into InvalidMatch struct to print out` <br> For example: <br>"error: number is above words number - (up, 999)" <br> "error: invalid number - hello (bin)" <br>  "error: invalid number (cap, hello) <br> 
switch case has a substring extracted by data[i][1] <br>
up|cap|low|bin|hex iterates from end of string and make modifications on words|digits based on number  <br> For example: <br> 
"Ready, set, go (up) !" -> "Ready, set, GO!" <br> "Welcome to the Brooklyn bridge (cap)" -> "Welcome to the Brooklyn Bridge" <br> "1E (hex) files were added" -> "30 files were added" <br> I should stop SHOUTING (low)" -> "I should stop shouting" <br>"It has been 10 (bin) years" -> "It has been 2 years" <br> "1E (hex) files were added" -> "30 files were added"
 <br>

- Block for punctuation formatting  <br>
**makePunct** formats punctuation, adding a space where required, every single or double quote should have closing ones, they should be placed to the right and left of the word without any spaces. Group punctuation should be close to each other with no space <br> For example: <br> "I was sitting over there ,and then BAMM !!" -> "I was sitting over there, and then BAMM!!" <br> '  I am in a good mood!    ' -> 'I am in a good mood!'

 - Block for article formatting <br>
Every a should be turned into an if the next word begins with a vowel (a, e, i, o, u) or a h. <br> For example: <br> "There it was. A amazing rock!" -> "There it was. An amazing rock!".