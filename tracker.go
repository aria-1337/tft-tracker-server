/* tracker.go
    This contains the bulk of the logic to handle calculation + storage of data analysis via the
    Riot instance from api.go. Namely methods that are not standard to the official api
*/
package main

type FullSummoner struct {
    RAccount RiotAccount
    RSummoner RiotSummoner
    MatchIds []string
    Matches []TFTMatch
}

func (r *RiotAPI) fetchSummonerByName(gameName string, tagLine string) FullSummoner {
    riotAccount := r.AccountByRiotId(gameName, tagLine)
    riotSummoner := r.SummonerByPuuid(riotAccount.Puuid)
    matchIds := r.MatchIdsByPuuid(riotAccount.Puuid)

    var matches []TFTMatch
    // TEST: only fetch 1 match
    for i := range matchIds {
        if i >= 1 {
            break
        }
        match := r.MatchById(matchIds[i])
        matches = append(matches, match)
    }

    return FullSummoner{
        RAccount: riotAccount,
        RSummoner: riotSummoner,
        MatchIds: matchIds,
        Matches: matches,
    }
}
