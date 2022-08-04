import {Lang}	from "@/components/graph/utils/lang";

export class Config {
	// -- DEFINE ---------
	global:{
		//graph:"topology"|"hierarchy"
		padding: { top:number, left:number, right:number, bottom:number }
		toolbar: {
			visible:boolean
			align: {
				horizontal:"none"|"left"|"right"|"center"
				vertical:"none"|"top"|"bottom"|"middle" 
			}
			margin: { left:number, top:number, right:number, bottom:number} 
		}
		scale: {
			ratio:number, maxRatio:number, minRatio:number
		}
	}
	data?:any
	extends: {
		hierarchy: {
			type: "horizontal"|"vertical"
			scale: {
				minWidth:number, maxWidth:number
			}
			align: {
				vertical:"none"|"top"|"bottom"|"middle"
			}
			group: {
				divide:boolean
				spacing:number
				title: {
					display:"always"|"none"|"has"
					spacing:number
				}
				box: {
					padding: { top:number, left:number, right:number, bottom:number }
					background: { fill:string, opacity:number}
					border: { width:number, color?:string, dash?:string }
					tree : {
						spacing:number
						node : {
							height:number
							padding: { top:number, left:number, right:number, bottom:number }
						}
						line : {
							caption: {
								align: "center"|"left"|"right"
								padding: { left:number, right:number }
								width: number
							}
							end:"none"|"arrow"
						}
					}
				}
			}
			node: {
				forEach: any
			}
		}
		topology: {
			tick: {
				skip:number
			}
			collision: {
				radius:number
			},
			simulation: {
				alphaDecay:number,
				onEnd: any
			}
		}
	}
	on?: {
		nodeclick?:(this: SVGElement, event: any, d: any) => void
	}
	//merge: (conf:Config)=> Config
	// -- 생성자 - Default 값 ---------
	constructor() {
		this.global = {
			//graph:"topology",
			toolbar: {
				visible: true,
				align: { horizontal:"right", vertical:"top" },
				margin: { top: 0, left: 0, right:0, bottom:0 }
			},
			padding: { top: 0, left: 0, right:0, bottom:0 },
			scale: { ratio: 1, minRatio: 0.1, maxRatio: 10 },
		};
		this.extends = {
			hierarchy: {
				type: "horizontal",
				scale: { minWidth: 0, maxWidth:0 },
				align: { vertical:"middle" },
				group: {
					divide: true,		//그룹으로 나누기
					spacing:25,			//group간 간격
					title: {
						display: "has",	//group title visible/hidden
						spacing: 10		//group title과 box 사이 간격
					},
					box: {
						border: { width: 1, color:"gray", dash: "2 2" },	//box border
						background: { fill:"silver", opacity:0.1 },			//box background
						padding: {top:10, left:5, right:5, bottom:10 },		//box padding
						tree : { 
							spacing:15,				//트리간 간격
							node : {
								height: 30, 	//노드 높이
								padding: { top:5, left:10, right:10, bottom:5 }	// only veritcal-graph
							},	
							line: { 
								caption: {
									align: "center",	// 라인 설명 정렬
									padding: { left:0, right:0 },	// 라인 설명 padding
									width: 0						// width fix (0:flexible, less 1: % )
								},
								end: "none"			// 라인 종료 모양
							}
						}
					},
				},
				node: {
					forEach: undefined	// 노드 데이터 foreach 함수
				}
			},
			topology: {
				tick: { skip:10 },
				collision: { radius:60 },
				simulation: { alphaDecay:0.006, onEnd: undefined }
			}
		};
		//this.merge = (conf:Config) => {
		//	return Lang.merge(this, conf)
		//}

	}
	//public	merge(conf:Config) {
	//	return Lang.merge(this, conf)
	//}

}
