# CloudWalk Technical Assessment

# Summary

- [Instructions](#instructions)
- [Observations](#observations)
- [Tasks](#tasks)
  - [Core Functionality](#core-funcionality)
  - [Optional](#optional)

---

The main details about the task are in the file [instructions.md](/instructions.md).
Below are some instructions andd observations about the task to help maintain what should be done and which things should be kept in mind. Interesting observations that were found from the log are highlighted with **Caveat**.

## Instructions

### Prerequisites

- Go installed on your machine.

### Setup

1. Clone this repository:
   ```bash
   git clone git@github.com:wscnd/cloudwalk-tests-software-engineer.git
   cd cloudwalk-tests-software-engineer
   ```
2. Run the application in Docker:

   ```bash
    make docker-all
   ```

   This command:

   - Creates a image
   - Runs a container with it which runs the app
   - Copies the json output to a local file
   - Deletes the container.

After running the application successfully, `match_data.json` will be copied from the docker container.

3. (Optional) Delete generated files:

   ```bash
    make clean
   ```

---

---

## Observations

**Input**: Quake Game log file [qgames.log](/logs/qgames.log), each line represents a log entry containing matches information.

**Output**: Parsed log grouped by match with json structure as follows:

```json
"game_1": {
"total_kills": 45,
"players": ["Dono da bola", "Isgalamido", "Zeh"],
"kills": {
  "Dono da bola": 5,
  "Isgalamido": 18,
  "Zeh": 20
  },
  "kills_by_means": {
    "MOD_SHOTGUN": 10,
    "MOD_RAILGUN": 2,
    "MOD_GAUNTLET": 1,
    ...
  }
}
```

### Approach

#### 1. Identify how the lines are formatted

Each line that we are interested appears to have the format of
`<timestamp> <EventType>: <Event Data…>`.

#### 2. Identify a Match boundary

##### 2.1 Characteristics

- Each new match starts with `InitGame` events.
- Game data is in between two `InitGame`.
- **Caveat 1**: Sometimes games are not ended with `ShutdownGame` events.

##### 2.2 Testing

- I manually detected 21 games, maybe assert first that the processing detected these games.
- Assert that the 21(?) matches start and end have the correct boundaries with the timestamp.

#### 3. Gather Data from Matches by Events

##### 3.1. ClientUserinfoChanged

- Means that a player changed something, sometimes they change their name, which is probably the most relevant change.
- Structure:

```
<timestamp> ClientUserinfoChanged: <User_ID> n\<User_Name>\t….some\other\things
```

- Example:

```
21:15 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
```

- We are interested in whatever is between `n\` and `\t` which is the player nickname associated with the change.
- **Caveat 2**: Sometimes a player change its name, identify by the `<User_ID>` and persist its previous match data.
- **Caveat 3**: A player can also have a nickname that is composed of more than one word, ex: _"My very nice nickname"_.

##### 3.2. Kill

- Means that there was a kill.
- Structure:

```
<timestamp> Kill: <Killer_ID> <Victim_ID> <Death_Cause_ID>: <Killer as string> killed <Victim as string> by <Death_Cause as string>
```

- Example:

```
21:42 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
```

- When `<world[Killer_ID=1022]>` kills a player, that player loses -1 kill score.
- Since `<world[Killer_ID=1022]>` is not a player, it should not appear in the list of players or in the dictionary of kills.,
- The counter `total_kills` includes player and world deaths.
- **Caveat 4**: Sometimes a player kills itself `(Killer_ID ==Victim_ID)`, the Kill shouldn't count and only the Death. Since this is a team game, some parsers can process the Death differently depending on how it handles the individual player KDA metric or team statistics.

## Tasks

### Core Funcionality

#### 1. Identify a Match boundary

- [x] Parse the log file to identify match boundaries using `InitGame` events.
- [x] Games not ending with `ShutdownGame` events.
- [x] Test 21 games.
- [x] Test 21 games start and end based on timestamp.

#### 2. Gather Data from Matches by Events

##### ClientUserinfoChanged events

- [x] Track `ClientUserinfoChanged` events to maintain accurate player names.
- [x] Player name changes during a match.
- [x] Multi-word player nicknames

##### Kill events

- [x] Player kills.
- [x] Deaths by `<world[Killer_ID=1022]>`.
- [x] Self-kills (do not count as kills).
- [x] Ensure `total_kills` includes all deaths (player and `<world[Killer_ID=1022]>`).
- [x] Accurate handling of self-kills and deaths.
- [x] Classify Kills by Death Cause.

#### 3. Group data

- [x] Group parsed data by match and output in the specified JSON structure.

#### 4. Output

- [x] Grouped information for each match
- [x] A player ranking
- [x] Report of deaths grouped by death cause for each match
- [x] Print to stdout

### Optional

#### 1. Processing Optimization

- [x] Reduce processing time using concurrency

#### 2. Data Modeling

- [ ] Review and refine entity models (e.g., Match, Player, Event)
- [ ] Implement business logic within entities
- [ ] Define clear terminology for entity operations
