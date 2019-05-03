# jot-frame
A golang library for writing dynamic content to the terminal

## Expansion
- no need for multiple frames? If I want to inject a frame within a frame that would be complicated. 
  However, if there was just one frame and you were manipulating lines, that works fine.
- Represent a graph data structure with live and dead elements
- Track lines with a line group, agnostic of the frame it sits within. This can help with "sub frames" with special purposes.
- I should be able to write to a line group as if it were a small terminal
- Certain lines should be pinned to the window, such that they never overflow from the window.
- need to rate limit writes to the screen as a whole
- hide/show lines

## Todo
- frame interface
- remove all
- add header, footer / remove header / footer
- turn on/off config options on the fly
- write convienence function? maybe not
- Size function that includes header and footer
- convienence line.Remove() to remove from the owning frame
- change Remove(*Line) to Remove(*Line...) ?
- Reopen line or frame (add Open())
