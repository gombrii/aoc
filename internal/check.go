package internal

func Check() error {
	//TODO:  Walk cache catalogue
	// Fo each lock == true create goroutine
	// in goroutine run puzzle.
	// Goroutines respond through a channel
	// 
	// While this is happening, on the main thread, one line is printed per puzzle/routine with a spinner
	// Each time a routine answers the corresponding line's spinner changes to a happy or sad symbol (maybe star or check mark). 


	return nil
}

//TODO: Also, concider not allowing locking of puzzles with other input than input.txt. It serves no purpose and confuses the conceptual model.
// What is locked is simply a puzzle (-d x -p x).

//TODO: Also, take chatGPT's advice and rearrange some of the commands. 
// Something like this:
//
// Usage:
//   aoc [puzzle run] -d DAY -p {1|2} [-y YEAR default: {{year}}] [-i INPUT default: input(.txt)]
//   aoc puzzle status  -d DAY -p {1|2} [-y YEAR] [-i INPUT]
//   aoc puzzle lock    -d DAY -p {1|2} [-y YEAR] [-i INPUT]
//   aoc puzzle unlock  -d DAY -p {1|2} [-y YEAR] [-i INPUT]
//   aoc init day       -d DAY [-y YEAR]
//   aoc init module    <name>
//   aoc cache clear
//   aoc help [command]
// 
// Commands:
//   puzzle run       Execute a puzzle (default when no subcommand is given)
//   puzzle status    Show last result and duration for a puzzle
//   puzzle lock      Freeze result; future runs must match it; keep best (fastest) duration
//   puzzle unlock    Unfreeze result; future runs update result and duration
//   init day         Scaffold solution files for a day
//   init module      Create a new AoC module structure
//   cache clear      Clear cached runners and metadata
// 
// Concepts:
//   • Puzzle = (year, day, part, input). Each puzzle has a cached runner plus:
//       res (last result) · dur (best/last duration*) · lock (true/false)
//   • Lock on → runs must match stored result; duration only updates if faster.
//   • Input: pass basename or filename. If no extension, “.txt” is appended.
//     Examples: -i test → test.txt,  -i data/example.txt → data/example.txt