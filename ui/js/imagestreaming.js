function ImageStreamer(imageSelecter,imageList,manualStartStop){
	var me=this
	var ws;
	var connected=false;
	manualStartStop.click(function(){
		manualStartStop.addClass('stream-change-wait')
		if(connected){
			me.stop()
		}else{
			me.start();
		}
	});
	me.onStatusMessage=function(n){ alert(n.message); }
	var imageClick=function(){
		imageSelecter.attr('src',$(this).attr('src'))
	}
	var processMessage = function (evt) 
	{ 
		var received_msg = evt.data;
		if( typeof(received_msg) == "string"){
			me.onStatusMessage(JSON.parse(received_msg));
			return;
		} 
		received_msg.type="image/jpeg";
	
		var li=$(document.createElement("li"));
		var img=$(document.createElement("img"));
		var imgUrl= window.URL.createObjectURL(received_msg); 
		
		li.append(img)
		
		imageList.append(li);
		
		imageSelecter.attr('src',imgUrl)
		img.attr('src',imgUrl);
		img.bind('load',function(){
			imageList.scrollLeft(imageList[0].scrollWidth);
		})
		img.click(imageClick)
		
		
	};
	me.start=function(){
		ws = new WebSocket("ws://"+location.host+"/imagestream");
		ws.onmessage=processMessage;
		ws.onclose=function(){
			connected=false;
			manualStartStop.addClass('stream-disconnected')
			manualStartStop.removeClass('stream-connected')
			manualStartStop.removeClass('stream-change-wait')
		}
		ws.onerror=ws.onclose;
		ws.onopen=function(){
			connected=true;
			manualStartStop.addClass('stream-connected')
			manualStartStop.removeClass('stream-disconnected')
			manualStartStop.removeClass('stream-change-wait')
		}
	}
	me.stop=function(){
		if( ws ) {
			ws.close()
		}
	}
	
}
