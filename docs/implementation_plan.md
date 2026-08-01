# Plan to Achieve 100% Test Parity in Anitomy-Go

After analyzing the failing output from the `anitomy/test/data.json` dataset, I have categorized the test failures into a few core issues. Here is the step-by-step implementation plan to resolve them and achieve full test parity.

## Proposed Changes

### 1. Implement Missing Parsers
- Currently, `ParseVideoResolution` in `parse_metadata2.go` is stubbed to return `nil`. We need to implement it. It should search for common resolution formats (e.g. `1280x720`, `1920x1080`, `720p`, `1080p`) using regex and extract them correctly, as `video_resolution` elements are completely missing from the outputs.
- Verify `ParseVolume` behaves correctly, as it was hastily implemented with regex boundary issues.

### 2. Refine Episode Parsing (`parse_episode.go`)
- **Release Version Formatting**: The tests expect `2` but the parser currently returns `v2`. The regex capture group extraction needs to be adjusted to omit the `v` or `V` prefix.
- **Dangling Numbers**: The parser currently fails to recognize `01` in patterns like `Title - 01`. Because `01` is not identified as the episode, the `Title` parser greedily consumes it as part of the title (resulting in a title like `"Princess Lover! - 01"` instead of `"Princess Lover!"`). We need to accurately port the `first_number` and `last_number` heuristics from C++ to correctly isolate these standalone episode numbers.

### 3. Refine Episode Title Parsing (`parse_episode_title.go`)
- The parser is failing to locate titles like `Tiger and Dragon` in `[TaigaSubs]_Toradora!_(2008)_-_01v2_-_Tiger_and_Dragon`. This is because the episode title span logic depends heavily on the Episode token being correctly identified first, and relies on correctly bounded enclosed brackets to stop parsing.
- We will fix the `find_episode_title` boundaries to correctly stop when it hits an identified token or a bracket.

### 4. Refine Title Extraction (`parse_title.go`)
- Once the episode and episode title bounds are correctly identified, the title extraction will naturally fix itself, as it just scans for the first free unenclosed range of tokens that hasn't been identified.

## Verification Plan

### Automated Tests
- I will execute `go test -v ./...` in the `anitomy-go` directory after each major component fix.
- My goal is to see the test failures drop to zero (or as close to zero as possible) for the official `data.json` test suite.

## User Review Required
No significant architectural changes are proposed. This is entirely an iterative debugging and porting accuracy task. Please review the plan above and click **Proceed** if you approve.
