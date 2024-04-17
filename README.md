
# Go Go Keyboard Cadet!

A typing game, based on an old Amiga game of the same name that I dimly remember.
This is an informal project developed together with my son to learn typing & game programming; things will be a bit messy :)

Music & sound effects on the dirtywave M8.

It's "playable" but the current typing "curriculum" is not very good; see the todo below.


= todo
 - [x] menu settings (speed, mission choice)
 - [x] make speed setting affect the right thing.
 - [x] adjust scores based on difficulty.
 - [ ] add mothership event (n-key rollover?)
     - [ ] show which keys need to be pressed?
     - [x] mothership sound effect.
     - [x] mothership graphics
     - [ ] mothership explode?
     - [ ] moship victory sound?
 - [x] add 'get ready' phase.
 - [x] add "approaching" sfx/riser. 
 - [x] animate nonwords
 - [x] hand graphics - home
 - [ ] hand graphics -- other keys
 - [ ] victory screen score & progression message (high score, personal best, etc)
 - [x] 'crashed' display & message.
 - [x] defeat music.  womp womp!
 - [ ] display name of level / help text.
 - [ ] High scores screen.
 - [ ] Redo "lesson" curriculum / better word choices. (Shifts & capitals)
 - [ ] Refactor how "profiles" work.
 - [ ] game options screen 
 - [ ] help screen 
 - [ ] practice mode.
 - [ ] targeting reticule graphic?
 - [x] victory fanfare
 - [x] persist scores
 - [ ] add reverb to sound fx
 - [x] opening music
 - [ ] figure out opening music loop point.
 - [x] better "shields" graphic.
 - [ ] Convert all the sounds to ogg. (some are mono?)
 - [ ] bug: riser sounds very quiet on cros?


= Scenes

1. Title Sequence
2. Menu
3. Profile Save/Load
4. Options
5. Mission Select
6. Training
7. Game
8. Mission Report

== Game context object

  - Options/Settings
  - Profile
    - name
    - missions completed & scores
  - Mission config/ game scene info
  - game scene results

nice pattern for audio?
https://github.com/hajimehoshi/go-inovation/blob/main/ino/internal/audio/audio.go

see also: https://pkg.go.dev/github.com/tinne26/edau#section-readme
and https://github.com/SolarLune/resound


= Missions, one idea?

  - f/j
  - d/k
  - s/l
  - a/;
  - df/jk
  - sdf/jkl
  - asdf/jkl;
  - gh + home
  - e/i + home
  - rtyu + home
  - qwe/iop + home


  original missions:
  - home row
  - heo
  - ti L-shift
  - rw.'
  - ng R-shift
  - ucp
  - yx,
  - mz=
  - bq?
  - v"
  - 0-9
  - !@$#%^&*()-+
  - word review
  - sentences
  - paragraphs
  speeds: fast, superfast, hyperfast

# Design for profiles/first-time-start

If no profile active at startup, default to instructions scene & next transition to profile management/enrollment s
