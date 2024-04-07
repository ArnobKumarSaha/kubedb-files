// sample persons
var persons = {
     "id01" : { "name": "hoge", "act_num": 50000}
     ,"id02" : { "name": "piyo", "act_num": 100000} 
     ,"id03" : { "name": "huga", "act_num": 50000} 
     ,"id04" : { "name": "mone", "act_num": 100000}
     ,"id05" : { "name": "tako", "act_num": 50000}
     ,"id06" : { "name": "nasu", "act_num": 100000}
     ,"id07" : { "name": "poke", "act_num": 50000}
     ,"id08" : { "name": "paca", "act_num": 100000}
     ,"id09" : { "name": "pipi", "act_num": 50000}
     ,"id10" : { "name": "zizi", "act_num": 100000}
    };
    
    // sample host names
    var hostNames = ['sv1', 'sv2', 'sv3'];
    var randHostName = function(){
     return hostNames[Math.floor( Math.random() * 3 )];
    }
    
    // sample point
    var randAddPoint = function(){
     var point = Math.floor( Math.random() * 100 ) + 1;
     var bonus = Math.floor( point * 0.05 );
     return { 'total': point + bonus, 'detail' : [point, bonus] };
    }
    
    // sample document template
    var docTemplate = {
     'info': {'uid': '', 'ts' : '', 'hn' : ''},
     'val' : {'total': 0, 'add':[0, 0]}
    };
    
    // timestamp
    var TimeStamp = function(baseTimestamp){
     this.baseTimestamp = baseTimestamp;
     this.baseTime = Date.parse(baseTimestamp);
     this.newTime = this.baseTime;
     this.date = new Date();
    };
    TimeStamp.prototype.tekitou = function(){
     var addSec = (Math.floor( Math.random() * 2000 ));
     this.newTime = this.newTime + addSec;
     this.date.setTime(this.newTime);
     var timeStamp = this.date.getFullYear() 
             + "-" 
             + ( ("0"+(this.date.getMonth() + 1)).slice(-2) ) 
             + "-" 
            + ("0"+ this.date.getDate() ).slice(-2)
            + " " 
            + ("0"+ this.date.getHours() ).slice(-2)
            + ":" 
            + ("0"+ this.date.getMinutes() ).slice(-2)
            + ":" 
            + ("0"+ this.date.getSeconds() ).slice(-2);
    
     return timeStamp;
    }
    
    // insert sample
    var baseTimeStamp = "2014-03-24 00:00:00";
    var bulkDocs = [];
    var docTemplateString = JSON.stringify(docTemplate);
    for ( var id in persons ){
     print("-->" + id);
    
     var total = 0;
     var ts = new TimeStamp(baseTimeStamp)
     for( var i=0; i<persons[id]["act_num"]; i++ ){
      // temporary document
      var tmpDoc = JSON.parse(docTemplateString);
      // info's property
      tmpDoc.info.uid = id;
      tmpDoc.info.ts = ts.tekitou();
      tmpDoc.info.hn = randHostName();
    
      // val's property
      var addPoint = randAddPoint();
      total = total + addPoint.total;
      tmpDoc.val.total = total;
      tmpDoc.val.add  = addPoint.detail;
      // push
      bulkDocs.push(tmpDoc);
      if( i!=0 && i % 10000 == 0 ) {
      print("\t" + id + ":" + i);
       //printjson(bulkDocs.length);
       db.kubedb.insert(bulkDocs);
       bulkDocs = [];
      }
     }
    
     db.kubedb.insert(bulkDocs);
     bulkDocs = [];
    }