<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0"
           xmlns:brid="http://www.opengis.net/citygml/bridge/2.0"
           xmlns:tun="http://www.opengis.net/citygml/tunnel/2.0"
           xmlns:tran="http://www.opengis.net/citygml/transportation/2.0">
  <cityObjectMember>
    <bldg:Building gml:id="B1">
      <bldg:measuredHeight>10.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <brid:Bridge gml:id="BR1"><brid:class>1000</brid:class></brid:Bridge>
  </cityObjectMember>
  <cityObjectMember>
    <tun:Tunnel gml:id="TN1"><tun:class>2000</tun:class></tun:Tunnel>
  </cityObjectMember>
  <cityObjectMember>
    <tran:Road gml:id="RD1"><tran:function>highway</tran:function></tran:Road>
  </cityObjectMember>
</CityModel>
