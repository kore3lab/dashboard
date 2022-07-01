import * as d3			from "d3";
import {Config}			from "@/components/graph/model/config.model";
import {Toolbar}		from "@/components/graph/toolbar";
import {Bounds, WH, UI}	from "@/components/graph/utils/ui";
import {Transform}		from "@/components/graph/utils/transform";
import {Lang}			from "@/components/graph/utils/lang";

/**
 * Graph 베이스 클래스
 */
export abstract class GraphBase {
	
	private m_config:Config = new Config();
	private m_container:HTMLElement;
	private m_on:any;

	public svg:d3.Selection<SVGSVGElement, any, SVGElement, any>;			//svg element
	public outlineWrapEl:d3.Selection<SVGGElement,any,SVGElement,any>;		//svg > g.outlineWrap
	public outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>;			//svg > g.outlineWrap > g.outline
	public toolbarEl:d3.Selection<SVGGElement,any,SVGElement,any>			//svg > g.toolbar
	public zoomBehavior:d3.ZoomBehavior<any,any>;							//zoom
	public resizeListener:any

	constructor(container?:string, conf?:Config) {
		if(container) this.container(container);
		if(conf) this.config(conf);
	}
	public container<T extends (GraphBase|HTMLElement)>(_?:string):T {
		return _ ? (this.m_container = d3.select<HTMLElement, any>(_).node()!, <T><unknown>this) : <T>this.m_container
	}
	public config<T extends (GraphBase|Config)>(_?:Config):T {
		return _ ? (this.m_config = new Config(), this.m_config = Lang.merge(this.m_config, _), <T><unknown>this) : <T>this.m_config;
	}
	public data(_?:any):GraphBase|any {
		return _ ? (this.m_config.data = _, this) : this.m_config.data;
	}
	public on(name?:string, func?:(this: SVGElement, event: any, d: any) => void):GraphBase|any {
		if (name && !this.m_on)  this.m_on = {};
		return name ? (this.m_on[name] = func, this) : this.m_on;
	}

	protected getBounds(): Bounds {
		let bounds:Bounds =  new Bounds(this.container<HTMLElement>());
		// width, height padding 반영
		const conf:Config = this.config<Config>();
		if (conf.global.padding.left + conf.global.padding.right != 0) {
			bounds.width -= (conf.global.padding.left + conf.global.padding.right);
			bounds.height -= (conf.global.padding.top + conf.global.padding.bottom);
		}
		return  bounds
	}

	/**
	 * 주어진 데이터를 기준으로 그래프를 랜더링한다.
	 *
	 * @param container container HTML element 
	 * @param config config(with data)
	 */
	public render():GraphBase {

		if(arguments.length==1) this.config(arguments[0]);
		else if(arguments.length==2) (this.container(arguments[0]), this.config(arguments[1]));

		if(!this.m_container) return this;
		const containerEl:d3.Selection<any, any, any, any> = d3.select(this.m_container);
		const conf:Config = this.config()
		if (this.on()) conf.on = Lang.merge(conf.on, this.on());

		// svg
		let svg:d3.Selection<SVGSVGElement, any, SVGElement, any> = containerEl.select<SVGSVGElement>("svg");
		if(svg.size() == 0) svg = containerEl.append("svg");
		svg.attr("preserveAspectRatio", "xMidYMid meet")

		if (conf.global.padding.top  != 0) svg.style("padding-top", conf.global.padding.top);
		if (conf.global.padding.bottom  != 0) svg.style("padding-bottom", conf.global.padding.bottom);
		if (conf.global.padding.left  != 0) svg.style("padding-left", conf.global.padding.left);
		if (conf.global.padding.right  != 0) svg.style("padding-right", conf.global.padding.right);
		this.svg = svg;

		// container init width, height
		const bounds:Bounds = this.getBounds();
		const initWH:WH = {width:bounds.width, height:bounds.height};

		// svg > g.outlineWrapEl  (resize 영역, svg 크기를 커버)
		this.outlineWrapEl = svg.select("g.outlineWrap")
		if(this.outlineWrapEl.size() == 0) {
			this.outlineWrapEl = svg.append("g").attr("class","outlineWrap");
			this.outlineWrapEl.append("rect").attr("class","background").attr("fill","transparent").attr("width",initWH.width).attr("height",initWH.height);
		}

		// svg > g.outlineWrapEl > g.outline (zoom 영역, 그래프 랜더링 영역)
		this.outlineEl = this.outlineWrapEl.select("g.outline");
		if(this.outlineEl.size() > 0) this.outlineEl.remove();
		this.outlineEl = this.outlineWrapEl.append("g").attr("class","outline");

		// populate 
		if(conf.data) this.populate(this.outlineEl, initWH, conf);

		// zoom
		this.outlineWrapEl.attr("transform", "translate(0 0) scale(1)");
		this.outlineWrapEl.on(".zoom", null);
		this.zoomBehavior = d3.zoom().on("zoom", (event)=> { 
			this.outlineEl.attr("transform", event.transform); 
		});
		this.outlineWrapEl.call(this.zoomBehavior.transform, d3.zoomIdentity.translate(0,0).scale(1));
		this.outlineWrapEl.call(this.zoomBehavior);		//binding zoom
		this.outlineWrapEl.on("dblclick.zoom", null);	//zoom 더블클릭 이벤트 drop (because event bubbling)
		this.zoomBehavior.scaleTo(this.outlineWrapEl, conf.global.scale.ratio)
		this.zoomBehavior.scaleExtent([conf.global.scale.minRatio,conf.global.scale.maxRatio]);	//min-max ratio

		// window resize event
		this.resizeListener = () => {
			const bounds:Bounds = this.getBounds();
			this.svg.attr("width", bounds.width).attr("height", bounds.height);
			const k:number = Math.round(Math.min(bounds.width/(initWH.width), bounds.height/(initWH.height))*100)/100;
			Transform.instance(this.outlineWrapEl.node()!).scale(k);
			UI.align(this.toolbarEl.node()!, conf.global.toolbar.align.horizontal, conf.global.toolbar.align.vertical);
		}
		if(!d3.select(window).on("resize.updatesvg")) d3.select(window).on("resize.updatesvg", this.resizeListener );

		// toolbar
		if (conf.global.toolbar.visible) {
			this.toolbarEl = Toolbar.render(this);
			UI.align(this.toolbarEl.node()!, conf.global.toolbar.align.horizontal, conf.global.toolbar.align.vertical);
		}

		this.resizeListener();
		return this;

	}

	/**
	 * zoom
	 * @param ratio 배율 (1보다 작으면 축소, 1보다 크면 확대)
	 */
	public zoom(k?:number):GraphBase {

		if (k) {
			this.zoomBehavior.scaleTo(this.outlineWrapEl, k);
		} else {
			const bounds:Bounds =  new Bounds(this.container<HTMLElement>());
			const outline:DOMRect =  this.outlineEl.node()!.getBBox();
			const k:number = Math.min(bounds.width/outline.width,bounds.height/outline.height); //ratio
			this.zoomBehavior.scaleTo(this.outlineWrapEl, k); //scale to
			this.zoomBehavior.translateTo(this.outlineWrapEl, outline.x + outline.width/2, outline.y + outline.height/2);	//center to
		}

		return this;
	}

	/**
	 * zoom 계속 증감 처리
	 * @param ratio 배율 (1보다 작으면 축소, 1보다 크면 확대, 현재 배율에 곱하기)
	 */
	public zoomRatio(ratio:number):GraphBase {
		let transform:d3.ZoomTransform = d3.zoomTransform(this.outlineEl.node()!);
		this.zoom(transform.k*ratio)
		return this;
	}


	protected abstract populate(outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>, wh:WH, conf:Config):void;

}
