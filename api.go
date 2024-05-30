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
