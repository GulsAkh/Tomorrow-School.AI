# ASCII ART

###  
The program receives a string argument and outputs in the a graphic representation using ASCII.

Main function processes input as argument, <br> generates origin hash and stores in a separate file to avoid modification of the file. <br> 
Embeds the file and splits it by a newline <br> Adding ascii characters into map's key as values <br> <br> handleAscii function iterates over input and if rune in the string is available in map it will print lines <br> Handles a newline in input as a real newline by replacing literal ('\n') into string ("\n")  <br> There are 3 banners file provided: standard.txt, shadow.txt, thinkertoy.txt  <br> Example1: <br> go run . "hello" | cat -e <br>  
```
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
``` 
Example2: <br> go run . "Hello\nThere" | cat -e 
```
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $

```