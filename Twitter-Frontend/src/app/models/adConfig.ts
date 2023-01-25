export class AdConfig {
    tweetID:   string = ""
	residence: string = ""
	gender:    string = ""
	ageFrom:   number = 0 
	ageTo:     number = 0

    AdConfig(tweetID: string, residence: string, gender: string, ageFrom: number, ageTo: number) {
        this.tweetID = tweetID
        this.residence = residence;
        this.gender = gender;
        this.ageFrom = ageFrom;
        this.ageTo = ageTo;
    }
}