declare namespace User{

	interface Info{
		id?:number
		name ?:string
		createTime?:string
		username?:string
		password?:string
		headImage?:string
		status?:number
		role?:number
		mail?:string
		token?:string
		oauthProvider?:string
		oauthId?:string
		publiccode?:string
		logined?:boolean
	}
}