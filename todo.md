# TODO

- [x] save a PNG
- [x] create a goroutine for each pixel
- [x] use a single channel for all pixel goroutines
- [x] create a sphere with intersect method
- [x] cast ray for each pixel into all objects
- [x] colour based on ~~angle or something~~
  - used distance and some scaling magic
- [ ] fix bug causing edge glow
- [ ] Make and use MVP matrices
- [ ] organize code
  - separate ratracing code from goroutine code somehow
- [ ] intersect cubes
- [ ] intersect triangles
- [ ] first hit sphere encapsulation
- [ ] create JSON structore to define world and materials etc.
- [ ] write JSON parser and Go structures to represent the world
- [ ] add point lights which affect colour
- [ ] investigate effectiveness of concurrency strategy
  - Is one goroutine per pixel ideal?
  - Maybe one goroutine per core is ideal? #cores - 1?

**... some more fundamental things ...**

- [ ] order rays in some way that allows a low-res render of the whole image to come in before any area gets more detail filled in
  - fill more than just pixel coming in, expecting it to be overwritten as more detail is filled in
  - maybe pixels at x, y \(\epsilon\) 2^n then 2^n-1, etc.  gradually filling in squares half the size of the previous level
- [ ] live OpenGl renderer
  - render to a texture, draw quad with that texture
- [ ] some sort of space partitioning
  - KD tree?
  - octree?
- [ ] shadows
- [ ] reflection
- [ ] refraction
- some kind of GI