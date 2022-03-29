import * as d3					from "d3";
import * as d3Zoom				from "d3-zoom";
import {ConfigModel, Bounds}	from "./model/models";
import {Transform}				from "./utils/transform";
import {UI,Lang}				from "./utils/lang";

/**
 * Graph 베이스 클래스
 */
export abstract class GraphBase {
	
	private _bounds:Bounds;
	private _outlineEl:SVGGElement;
	private _graphEl:d3.Selection<SVGGElement,any,SVGElement,any>;
	private _svg:d3.Selection<SVGSVGElement, any, SVGElement, any>;
	private _zoom:d3Zoom.ZoomBehavior<Element,any>;
	private _container:HTMLElement;
	private _config:ConfigModel.Config = new ConfigModel.Config();
	private _beforeAlign:{vertical:string, horizonal:string} = {vertical:"",horizonal:"" };

	constructor(el:HTMLElement, conf?:ConfigModel.Config) {
		this._container = el;
		Object.assign(this._config, conf);
	}

	public config(_?:ConfigModel.Config):GraphBase|ConfigModel.Config {
		return arguments.length ? (this._config = Lang.merge(new ConfigModel.Config(), _), this) : this._config;
	}
	public data(_?:any):GraphBase|any {
		return arguments.length ? (this._config.data = _, this) : this._config.data;
	}
	public bounds(_?:Bounds):any {
		return arguments.length ? (this._bounds = _, this): this._bounds;
	}
	protected outlineEl(_?:SVGGElement):any {
		return arguments.length ? (this._outlineEl = _, this): this._outlineEl;
	}
	protected graphEl(_?:d3.Selection<SVGGElement,any,SVGElement,any>):any {
		return arguments.length ? (this._graphEl = _, this): this._graphEl;
	}
	public svg(_?:d3.Selection<SVGSVGElement, any, SVGElement, any>):any {
		return arguments.length ? (this._svg = _, this): this._svg;
	}
	public container(_?:HTMLElement):any {
		return arguments.length ? (this._container = _, this): this._container;
	}
	public zoomBehavior(_?:d3Zoom.ZoomBehavior<Element,any>):any {
		return arguments.length ? (this._zoom = _, this): this._zoom;
	}
	

	/**
	 * 주어진 데이터를 기준으로 그래프를 랜더링한다.
	 * @param data 
	 */
	public render():GraphBase {

		if(arguments.length==1) this.config(arguments[0]);
		else if(arguments.length==2) (this._container = arguments[0], this.config(arguments[1]));

		if(!this._container) return;
		let container:d3.Selection<any, any, any, any> = d3.select(this._container);
		
		// svg
		let svg:d3.Selection<SVGSVGElement, any, SVGElement, any> = container.select<SVGSVGElement>("svg");
		if(svg.size() == 0) svg = container.append("svg");

		//bound 계산, padding 반영
		let bounds:Bounds =  new Bounds(container);

		// svg 크기 지정
		svg.attr("width", bounds.width).attr("height", bounds.height);


		let graphEl:d3.Selection<SVGGElement,any,SVGElement,any> = svg.select("g.graph");
		let outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>;
		let attrTransform:string;

		if(graphEl.size() > 0) {
			// 이전에 outline 이 있다면  
			// 		- 이전에 정렬정보가 변경되지 않고
			//		- transform 속성값을 삭제하기 전에 저장하여 나중에 재 설정해준다.
			outlineEl = svg.select<SVGGElement>("g.outline");
			if(outlineEl.size() > 0) {
				attrTransform = (this._beforeAlign.vertical  != this._config.global.align.vertical ||  this._beforeAlign.horizonal  != this._config.global.align.horizonal ) ? "": outlineEl.attr("transform")
				outlineEl.remove();
			}
			this._beforeAlign  = Object.assign({}, this._config.global.align);

		} else {
			// g.graph > g.outline 추가 
			//		- 그래프는 g.outline 에 추가됨
			//		- g.graph 는 zoom 영역임, svg 크기를 커버
			graphEl = svg.append("g").attr("class","graph");
			graphEl.append("rect").attr("class","background").attr("width",bounds.width).attr("height",bounds.height).attr("fill","transparent")
		}

		outlineEl = graphEl.append("g").attr("class","outline");


		// 멤버변수들
		this.bounds(bounds);
		this.outlineEl(outlineEl.node());
		this.svg(svg);
		this.container(container.node());
		this.graphEl(graphEl);

		// 데이터 모델 구성 후 그려주기
		this.populate(this._config, svg, bounds, outlineEl);

		// 이전에 outline 이 있다면  이전 속성값 다시 지정하고 수직, 수평정렬은 수행하지 않는다.
		if(attrTransform) outlineEl.attr("transform",attrTransform);
		else {

			// outline 수직 정렬
			if(this._config.global.align.vertical=="middle") UI.alignVertical(outlineEl.node());
			else if(this._config.global.align.vertical=="none") Transform.instance(outlineEl.node()).shiftY(this._config.global.padding.top);	//중간 정렬이 아닐 경우 TOP padding 적용

			// outline 수평 정렬
			if(this._config.global.align.horizonal=="center") UI.alignHorizonal(outlineEl.node());
			else if(this._config.global.align.horizonal=="none") Transform.instance(outlineEl.node()).shiftX(this._config.global.padding.left);	//가운데 정렬이 아닐 경우 LEFT padding 적용
		}

		// ZOOM
		this.zoomBehavior(
			d3Zoom.zoom().on("zoom", (event)=> {
				outlineEl.attr("transform", event.transform);  
			})
		);

		// ZOOM 가운데 정렬에 따를 초기화
		let transform:Transform = Transform.instance(this.outlineEl());
		graphEl.call(this.zoomBehavior().transform, d3Zoom.zoomIdentity.translate(transform.x, transform.y).scale(transform.k));

		// ZOOM 바인딩
		graphEl.call(this.zoomBehavior());
		graphEl.on("dblclick.zoom", null);	//zoom 더블클릭 이벤트 drop (because event bubbling)

		return this;

	}

	/**
	 * ZOOM 계속 증감 처리
	 * 
	 * @param ratio 배율 (1보다 작으면 축소, 1보다 크면 확대)
	 */
	public zoomRatio(ratio:number,range?:Array<number>):GraphBase {

		let transform = d3Zoom.zoomTransform(this.graphEl().node());
		let k:number = transform.k*ratio;
		if(range) {
			if(transform.k < range[0])  k = range[0];	//최소배율
			if(transform.k > range[1])  k = range[1];	//최대배율
		}

		this.graphEl().call(this.zoomBehavior().transform, d3Zoom.zoomIdentity.translate(transform.x, transform.y).scale(k));

		return this;
	}


	/**
	 * ZOOM
	 * 
	 * @param k 배율 (0 이면 SVG에 맞춤)
	 */
	public zoom(k?:number):GraphBase {

		if(k) {
			let transform = d3Zoom.zoomTransform(this.graphEl().node());
			this.graphEl().call(this.zoomBehavior().transform, d3Zoom.zoomIdentity.translate(transform.x, transform.y).scale(k));
		} else {

			Transform.instance(this.outlineEl()).translate(0,0).ratioScale(1);	//초기화

			let rect = this.outlineEl().getBoundingClientRect();
			let bounds:DOMRect =  this.bounds();
			
			let transform:Transform = new Transform(this.outlineEl());

			rect.width = rect.width * 1/transform.k;
			rect.height = rect.height * 1/transform.k;
			transform.k = 1;

			// k 결정 (너비, 높이)
			if(rect.width>bounds.width) transform.k = bounds.width/rect.width;
			if(rect.height>bounds.height && transform.k>bounds.height/rect.height) transform.k = bounds.height/rect.height;
			
			// x,y 축 이동 (가운데 정렬)
			transform.x = (bounds.width-(rect.width*transform.k))/2-rect.x;
			transform.y = (bounds.height-(rect.height* transform.k))/2-rect.y;

			this.graphEl().call(this.zoomBehavior().transform, d3Zoom.zoomIdentity.translate(transform.x, transform.y).scale(transform.k==0?1:transform.k));
		}

		return this;
	}

	/**
	 * 리사이즈 처리
	 * 
	 */
	public resize() {

		if(!this.bounds()) return;
		
		let w:number = this.bounds().width;
		let h:number = this.bounds().height;
		let transform:Transform = Transform.instance(this._outlineEl);

		let b:ClientRect = this._container.getBoundingClientRect();
		this.svg().attr("width", b.width).attr("height", b.height);
		this.svg().select("g.graph").select("rect.background").attr("width", b.width).attr("height", b.height);	// g.graph > rect.background 엘리먼트 크기도 조정

		w = (b.width == w ? 1: b.width/ w);
		h = (b.height== h ? 1: b.height/ h);
		let k:number = Math.min(w,h);
		if(k==1) return;

		transform.x = transform.x* k
		transform.y = transform.y* k
		transform.k = transform.k* k

		this.graphEl().call(this.zoomBehavior().transform, d3Zoom.zoomIdentity.translate(transform.x, transform.y).scale(transform.k==0?1:transform.k));
		this.bounds(<DOMRect>this._container.getBoundingClientRect());

	}
	
	protected abstract populate(conf:ConfigModel.Config, svg:d3.Selection<SVGSVGElement, any, SVGElement, any>, bounds:Bounds, outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>):void;

}
