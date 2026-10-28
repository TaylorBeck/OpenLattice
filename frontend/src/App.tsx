import React, { useEffect, useState } from 'react';
import { Viewer, Entity } from 'resium';
import { Cartesian3, Color } from 'cesium';
// import { openlattice } from './pb/track'; // Protobuf generated types

export default function App() {
  const [wsStatus, setWsStatus] = useState<'DISCONNECTED' | 'CONNECTED'>('DISCONNECTED');

  useEffect(() => {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/stream';
    const socket = new WebSocket(wsUrl);
    socket.binaryType = 'arraybuffer';

    socket.onopen = () => setWsStatus('CONNECTED');
    socket.onclose = () => setWsStatus('DISCONNECTED');

    socket.onmessage = (event) => {
      // TODO: Implement Protobuf decoding here
      // const streamData = openlattice.telemetry.TrackUpdateStream.decode(new Uint8Array(event.data));
    };

    return () => socket.close();
  }, []);

  return (
    <div style={{ width: '100vw', height: '100vh', margin: 0, overflow: 'hidden' }}>
      <div style={{
        position: 'absolute', top: 10, left: 10, zIndex: 100,
        color: '#00ff00', fontFamily: 'monospace', background: 'rgba(0,0,0,0.7)', padding: '8px'
      }}>
        LATTICE C2 LINK: {wsStatus}
      </div>

      <Viewer animation="{false}" baseLayerPicker="{false}" full geocoder="{false}" navigationHelpButton="{false}" timeline="{false}">
        {/* Entities will be dynamically mapped here based on the Protobuf stream */}
      </Viewer>
    </div>
  );
}