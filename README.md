# pwr_chess

Chess game with some additions made to learn golang.

# Features

- [x] Local multiplayer board.
- [ ] Room-based multiplayer.
- [x] SSR.
- [ ] API.
- [x] 0-indexed cell notation to make moves.
- [ ] Mathematical notation to make moves.
- [ ] Move validation.
  - [ ] Pawn.
    - [ ] General.
    - [ ] en passant.
    - [ ] Transformation.
  - [ ] King.
    - [ ] General.
    - [ ] No move to checked field.
    - [ ] Check.
    - [ ] Checkmate.
    - [ ] Castle.
  - [ ] Queen.
  - [ ] Rook.
  - [ ] Bishop.
  - [x] Knight.
- [ ] Piece kill count.
- [ ] Rule sets.
  - [ ] Classic.
  - [ ] Score-based (kill count).
  - [ ] Instagib (check == checkmate)
- [ ] Game replay.
- [ ] Undo tree.

# How to install

```bash
git clone https://github.com/F1encko627/pwr_chess.git
cd pwr_chess
go run main.go
```
