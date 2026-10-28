import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace openlattice. */
export namespace openlattice {

    /** Namespace telemetry. */
    namespace telemetry {

        /** Properties of a Track. */
        interface ITrack {

            /** Track entityId */
            entityId?: (string|null);

            /** Track class */
            "class"?: (openlattice.telemetry.Track.TrackClass|null);

            /** Track latitude */
            latitude?: (number|null);

            /** Track longitude */
            longitude?: (number|null);

            /** Track altitudeMeters */
            altitudeMeters?: (number|null);

            /** Track velocityKnots */
            velocityKnots?: (number|null);

            /** Track headingDegrees */
            headingDegrees?: (number|null);

            /** Track timestampMs */
            timestampMs?: (number|Long|null);
        }

        /** Represents a Track. */
        class Track implements ITrack {

            /**
             * Constructs a new Track.
             * @param [properties] Properties to set
             */
            constructor(properties?: openlattice.telemetry.ITrack);

            /** Track entityId. */
            public entityId: string;

            /** Track class. */
            public class: openlattice.telemetry.Track.TrackClass;

            /** Track latitude. */
            public latitude: number;

            /** Track longitude. */
            public longitude: number;

            /** Track altitudeMeters. */
            public altitudeMeters: number;

            /** Track velocityKnots. */
            public velocityKnots: number;

            /** Track headingDegrees. */
            public headingDegrees: number;

            /** Track timestampMs. */
            public timestampMs: (number|Long);

            /**
             * Creates a new Track instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Track instance
             */
            public static create(properties?: openlattice.telemetry.ITrack): openlattice.telemetry.Track;

            /**
             * Encodes the specified Track message. Does not implicitly {@link openlattice.telemetry.Track.verify|verify} messages.
             * @param message Track message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: openlattice.telemetry.ITrack, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Track message, length delimited. Does not implicitly {@link openlattice.telemetry.Track.verify|verify} messages.
             * @param message Track message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: openlattice.telemetry.ITrack, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Track message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Track
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): openlattice.telemetry.Track;

            /**
             * Decodes a Track message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Track
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): openlattice.telemetry.Track;

            /**
             * Verifies a Track message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Track message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Track
             */
            public static fromObject(object: { [k: string]: any }): openlattice.telemetry.Track;

            /**
             * Creates a plain object from a Track message. Also converts values to other types if specified.
             * @param message Track
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: openlattice.telemetry.Track, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Track to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for Track
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        namespace Track {

            /** TrackClass enum. */
            enum TrackClass {
                UNKNOWN = 0,
                CIVILIAN_AIR = 1,
                CIVILIAN_MARITIME = 2,
                THREAT_AIR = 3,
                THREAT_MARITIME = 4
            }
        }

        /** Properties of a TrackUpdateStream. */
        interface ITrackUpdateStream {

            /** TrackUpdateStream tracks */
            tracks?: (openlattice.telemetry.ITrack[]|null);
        }

        /** Represents a TrackUpdateStream. */
        class TrackUpdateStream implements ITrackUpdateStream {

            /**
             * Constructs a new TrackUpdateStream.
             * @param [properties] Properties to set
             */
            constructor(properties?: openlattice.telemetry.ITrackUpdateStream);

            /** TrackUpdateStream tracks. */
            public tracks: openlattice.telemetry.ITrack[];

            /**
             * Creates a new TrackUpdateStream instance using the specified properties.
             * @param [properties] Properties to set
             * @returns TrackUpdateStream instance
             */
            public static create(properties?: openlattice.telemetry.ITrackUpdateStream): openlattice.telemetry.TrackUpdateStream;

            /**
             * Encodes the specified TrackUpdateStream message. Does not implicitly {@link openlattice.telemetry.TrackUpdateStream.verify|verify} messages.
             * @param message TrackUpdateStream message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: openlattice.telemetry.ITrackUpdateStream, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified TrackUpdateStream message, length delimited. Does not implicitly {@link openlattice.telemetry.TrackUpdateStream.verify|verify} messages.
             * @param message TrackUpdateStream message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: openlattice.telemetry.ITrackUpdateStream, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a TrackUpdateStream message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns TrackUpdateStream
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): openlattice.telemetry.TrackUpdateStream;

            /**
             * Decodes a TrackUpdateStream message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns TrackUpdateStream
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): openlattice.telemetry.TrackUpdateStream;

            /**
             * Verifies a TrackUpdateStream message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a TrackUpdateStream message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns TrackUpdateStream
             */
            public static fromObject(object: { [k: string]: any }): openlattice.telemetry.TrackUpdateStream;

            /**
             * Creates a plain object from a TrackUpdateStream message. Also converts values to other types if specified.
             * @param message TrackUpdateStream
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: openlattice.telemetry.TrackUpdateStream, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this TrackUpdateStream to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for TrackUpdateStream
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        /** Properties of an Ack. */
        interface IAck {

            /** Ack success */
            success?: (boolean|null);

            /** Ack message */
            message?: (string|null);
        }

        /** Represents an Ack. */
        class Ack implements IAck {

            /**
             * Constructs a new Ack.
             * @param [properties] Properties to set
             */
            constructor(properties?: openlattice.telemetry.IAck);

            /** Ack success. */
            public success: boolean;

            /** Ack message. */
            public message: string;

            /**
             * Creates a new Ack instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Ack instance
             */
            public static create(properties?: openlattice.telemetry.IAck): openlattice.telemetry.Ack;

            /**
             * Encodes the specified Ack message. Does not implicitly {@link openlattice.telemetry.Ack.verify|verify} messages.
             * @param message Ack message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: openlattice.telemetry.IAck, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Ack message, length delimited. Does not implicitly {@link openlattice.telemetry.Ack.verify|verify} messages.
             * @param message Ack message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: openlattice.telemetry.IAck, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes an Ack message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Ack
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): openlattice.telemetry.Ack;

            /**
             * Decodes an Ack message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Ack
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): openlattice.telemetry.Ack;

            /**
             * Verifies an Ack message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates an Ack message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Ack
             */
            public static fromObject(object: { [k: string]: any }): openlattice.telemetry.Ack;

            /**
             * Creates a plain object from an Ack message. Also converts values to other types if specified.
             * @param message Ack
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: openlattice.telemetry.Ack, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Ack to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for Ack
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        /** Represents a ThreatService */
        class ThreatService extends $protobuf.rpc.Service {

            /**
             * Constructs a new ThreatService service.
             * @param rpcImpl RPC implementation
             * @param [requestDelimited=false] Whether requests are length-delimited
             * @param [responseDelimited=false] Whether responses are length-delimited
             */
            constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

            /**
             * Creates new ThreatService service using the specified rpc implementation.
             * @param rpcImpl RPC implementation
             * @param [requestDelimited=false] Whether requests are length-delimited
             * @param [responseDelimited=false] Whether responses are length-delimited
             * @returns RPC service. Useful where requests and/or responses are streamed.
             */
            public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): ThreatService;

            /**
             * Calls InjectTrack.
             * @param request Track message or plain object
             * @param callback Node-style callback called with the error, if any, and Ack
             */
            public injectTrack(request: openlattice.telemetry.ITrack, callback: openlattice.telemetry.ThreatService.InjectTrackCallback): void;

            /**
             * Calls InjectTrack.
             * @param request Track message or plain object
             * @returns Promise
             */
            public injectTrack(request: openlattice.telemetry.ITrack): Promise<openlattice.telemetry.Ack>;
        }

        namespace ThreatService {

            /**
             * Callback as used by {@link openlattice.telemetry.ThreatService#injectTrack}.
             * @param error Error, if any
             * @param [response] Ack
             */
            type InjectTrackCallback = (error: (Error|null), response?: openlattice.telemetry.Ack) => void;
        }
    }
}
