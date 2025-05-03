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
    - [x] General movement.
    - [x] Taking pieces.
    - [ ] En passant.
    - [ ] Transformation.
  - [ ] King.
    - [x] General.
    - [x] No move to checked field.
    - [x] Check.
    - [ ] Checkmate (REALLY HARD).
    - [ ] Castle.
  - [x] Queen.
  - [x] Rook.
  - [x] Bishop.
  - [x] Knight.
- [ ] Piece kill count.
- [ ] Game modes.
  - [ ] Classic.
  - [ ] Score-based (kill count).
  - [ ] Instagib (check == checkmate)
  - [ ] Blitz (timer)
- [ ] Game replay.
- [ ] Undo tree.
- [ ] Better interface.

# How to install

```bash
git clone https://github.com/F1encko627/pwr_chess.git
cd pwr_chess
go run main.go
```

# TODO INSIGHTES

- For game modes... store pointer of MovePiece() function or separate move validation functions.
- For undo tree store diff (just write down consequenses of current move).
- Add new optional parameter for moving piece to transform pawn. Ignored pawn not on the last row.
- Store pawns ability to en passant for every board, clear after move.
- After check, check for mate:
  1. king can't move over in any way (move or take any piece);
  2. any piece cant take (check every piece abilty to move from current position to piece giving check, using already written move validation functions);
  3. any piece can't shield (HARDEST);
