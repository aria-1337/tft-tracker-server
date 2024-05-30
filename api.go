package main

import (
    "encoding/json"
    "net/http"
    "io/ioutil"
    "fmt"
)

// Interface to official RIOT API
type RiotAPI struct {
    Token string
    Region string
    BaseUrl string
    AlternateUrl string
}

func Riot(token string, region string) (*RiotAPI) {
    baseUrl, alternateUrl := getUrlPrefixs(region)

    riot := &RiotAPI{
        Token: token,
        Region: region,
        BaseUrl: baseUrl,
        AlternateUrl: alternateUrl,
    }

    return riot
}

func getUrlPrefixs(baseRegion string) (string, string) {
    regionMap := map[string]string{
        "na1": "americas",
    }

    baseUrl := "https://" + baseRegion + ".api.riotgames.com/"
    alternateUrl := "https://" + regionMap[baseRegion] + ".api.riotgames.com/"
    
    return baseUrl, alternateUrl
}

// Generic request, still need to decode body
func (r *RiotAPI) Request(url string, target interface{}) {
    fmt.Println("New request:", url)

    client := http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    
    if err != nil {
        fmt.Println(err)
    }

    req.Header = http.Header{
        "X-Riot-Token": {r.Token},
    }

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }

    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if err := json.Unmarshal(body, &target); err != nil {
        fmt.Println("Error decoding JSON in AccountByRiotId")
    }
}

type RiotAccount struct {
    Puuid string `json:"puuid"`
    GameName string `json:"gameName"`
    TagLine string `json:"tagLine"`
}

func (r *RiotAPI) AccountByRiotId(gameName string, tagLine string) RiotAccount {
    url := r.AlternateUrl + "riot/account/v1/accounts/by-riot-id/" + gameName + "/" + tagLine
    var result RiotAccount 
    r.Request(url, &result)

    return result
}

type RiotSummoner struct {
    Id string `json:"id"`
    AccountId string `json:"accountId"`
    Puuid string `json:"puuid"`
    ProfileIconId int `json:"profileIconId"`
    SummonerLevel int `json:"summonerLevel"`
}

func (r *RiotAPI) SummonerByPuuid(puuid string) RiotSummoner {
    url := r.BaseUrl + "tft/summoner/v1/summoners/by-puuid/" + puuid
    var result RiotSummoner
    r.Request(url, &result)

    return result
}

func (r *RiotAPI) MatchIdsByPuuid(puuid string) []string {
    // TODO: Can add count start endtime starttime
    url := r.AlternateUrl + "tft/match/v1/matches/by-puuid/" + puuid + "/ids?start=0&count=50"
    var result []interface{}
    r.Request(url, &result)

    var cleanIds []string
    for i := range result {
        cleanIds = append(cleanIds, result[i].(string))
    }

    return cleanIds
}

type TFTMatch struct {
    Metadata struct {
        DataVersion string `json:"data_version"`
        MatchId string `json:"match_id"`
        Participants []string `json:"participants"`
    }
    Info struct {
        EndOfGameResult string `json:"endOfGameResult"`
        GameCreation int `json:"gameCreation"`
        GameId int `json:"gameId"`
        GameDateTime int `json:"game_datetime"`
        GameLength int `json:"game_length"`
        GameVersion string `json:"game_version"`
        MapId int `json:"mapId"`
        Participants []struct {
            Augments []string `json:"augments"`
            Companion struct {
                ContentId string `json:"content_ID"`
                ItemId int `json:"item_ID"`
                SkinId int `json:"skin_ID"`
                Species string `json:"species"`
            }
            GoldLeft int `json:"gold_left"`
            LastRound int `json:"last_round"`
            Level int `json:"level"`
            Placement int `json:"placement"`
            Puuid string `json:"puuid"`
            TimeEliminated int `json:"time_eliminated"`
            TotalDamageToPlayers int `json:"total_damage_to_players"`
            Traits []struct {
                Name string `json:"name"`
                NumUnits int `json:"num_units"`
                TierCurrent int `json:"tier_current"`
                TierTotal int `json:"tier_total"`
            }
            Units []struct {
                CharacterId string `json:"character_id"`
                ItemNames []string `json:"itemNames"`
                Name string `json:"name"`
                Rarity int `json:"rarity"`
                Tier int `json:"tier"`
            }
        }
        QueueId int `json:"queueId"`
        GameType string `json:"tft_game_type"`
        SetName string `json:"tft_set_core_name"`
        SetNumber int `json:"tft_set_number"`
    }
}

func (r *RiotAPI) MatchById(matchId string) TFTMatch {
    url := r.AlternateUrl + "tft/match/v1/matches/" + matchId
    var result TFTMatch
    r.Request(url, &result)

    return result
}
