/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const openlattice = $root.openlattice = (() => {

    /**
     * Namespace openlattice.
     * @exports openlattice
     * @namespace
     */
    const openlattice = {};

    openlattice.telemetry = (function() {

        /**
         * Namespace telemetry.
         * @memberof openlattice
         * @namespace
         */
        const telemetry = {};

        telemetry.Track = (function() {

            /**
             * Properties of a Track.
             * @memberof openlattice.telemetry
             * @interface ITrack
             * @property {string|null} [entityId] Track entityId
             * @property {openlattice.telemetry.Track.TrackClass|null} ["class"] Track class
             * @property {number|null} [latitude] Track latitude
             * @property {number|null} [longitude] Track longitude
             * @property {number|null} [altitudeMeters] Track altitudeMeters
             * @property {number|null} [velocityKnots] Track velocityKnots
             * @property {number|null} [headingDegrees] Track headingDegrees
             * @property {number|Long|null} [timestampMs] Track timestampMs
             */

            /**
             * Constructs a new Track.
             * @memberof openlattice.telemetry
             * @classdesc Represents a Track.
             * @implements ITrack
             * @constructor
             * @param {openlattice.telemetry.ITrack=} [properties] Properties to set
             */
            function Track(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null && keys[i] !== "__proto__")
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Track entityId.
             * @member {string} entityId
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.entityId = "";

            /**
             * Track class.
             * @member {openlattice.telemetry.Track.TrackClass} class
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype["class"] = 0;

            /**
             * Track latitude.
             * @member {number} latitude
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.latitude = 0;

            /**
             * Track longitude.
             * @member {number} longitude
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.longitude = 0;

            /**
             * Track altitudeMeters.
             * @member {number} altitudeMeters
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.altitudeMeters = 0;

            /**
             * Track velocityKnots.
             * @member {number} velocityKnots
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.velocityKnots = 0;

            /**
             * Track headingDegrees.
             * @member {number} headingDegrees
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.headingDegrees = 0;

            /**
             * Track timestampMs.
             * @member {number|Long} timestampMs
             * @memberof openlattice.telemetry.Track
             * @instance
             */
            Track.prototype.timestampMs = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

            /**
             * Creates a new Track instance using the specified properties.
             * @function create
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {openlattice.telemetry.ITrack=} [properties] Properties to set
             * @returns {openlattice.telemetry.Track} Track instance
             */
            Track.create = function create(properties) {
                return new Track(properties);
            };

            /**
             * Encodes the specified Track message. Does not implicitly {@link openlattice.telemetry.Track.verify|verify} messages.
             * @function encode
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {openlattice.telemetry.ITrack} message Track message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Track.encode = function encode(message, writer, q) {
                if (!writer)
                    writer = $Writer.create();
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                if (message.entityId != null && Object.hasOwnProperty.call(message, "entityId"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.entityId);
                if (message["class"] != null && Object.hasOwnProperty.call(message, "class"))
                    writer.uint32(/* id 2, wireType 0 =*/16).int32(message["class"]);
                if (message.latitude != null && Object.hasOwnProperty.call(message, "latitude"))
                    writer.uint32(/* id 3, wireType 1 =*/25).double(message.latitude);
                if (message.longitude != null && Object.hasOwnProperty.call(message, "longitude"))
                    writer.uint32(/* id 4, wireType 1 =*/33).double(message.longitude);
                if (message.altitudeMeters != null && Object.hasOwnProperty.call(message, "altitudeMeters"))
                    writer.uint32(/* id 5, wireType 1 =*/41).double(message.altitudeMeters);
                if (message.velocityKnots != null && Object.hasOwnProperty.call(message, "velocityKnots"))
                    writer.uint32(/* id 6, wireType 1 =*/49).double(message.velocityKnots);
                if (message.headingDegrees != null && Object.hasOwnProperty.call(message, "headingDegrees"))
                    writer.uint32(/* id 7, wireType 1 =*/57).double(message.headingDegrees);
                if (message.timestampMs != null && Object.hasOwnProperty.call(message, "timestampMs"))
                    writer.uint32(/* id 8, wireType 0 =*/64).int64(message.timestampMs);
                return writer;
            };

            /**
             * Encodes the specified Track message, length delimited. Does not implicitly {@link openlattice.telemetry.Track.verify|verify} messages.
             * @function encodeDelimited
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {openlattice.telemetry.ITrack} message Track message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Track.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer && writer.len ? writer.fork() : writer).ldelim();
            };

            /**
             * Decodes a Track message from the specified reader or buffer.
             * @function decode
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {openlattice.telemetry.Track} Track
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Track.decode = function decode(reader, length, error, long) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                if (long === undefined)
                    long = 0;
                if (long > $Reader.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let end, message;
                if (length === undefined)
                    end = reader.len;
                else {
                    end = reader.pos + length;
                    if (end > reader.len)
                        throw RangeError("index out of range");
                    length = reader.len;
                    reader.len = end;
                }
                message = new $root.openlattice.telemetry.Track();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    if (tag === error)
                        break;
                    switch (tag >>> 3) {
                    case 1: {
                            message.entityId = reader.string();
                            break;
                        }
                    case 2: {
                            message["class"] = reader.int32();
                            break;
                        }
                    case 3: {
                            message.latitude = reader.double();
                            break;
                        }
                    case 4: {
                            message.longitude = reader.double();
                            break;
                        }
                    case 5: {
                            message.altitudeMeters = reader.double();
                            break;
                        }
                    case 6: {
                            message.velocityKnots = reader.double();
                            break;
                        }
                    case 7: {
                            message.headingDegrees = reader.double();
                            break;
                        }
                    case 8: {
                            message.timestampMs = reader.int64();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7, long);
                        break;
                    }
                }
                if (length !== undefined) {
                    if (reader.pos !== end)
                        throw RangeError("index out of range");
                    reader.len = length;
                }
                return message;
            };

            /**
             * Decodes a Track message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {openlattice.telemetry.Track} Track
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Track.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Track message.
             * @function verify
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Track.verify = function verify(message, long) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    return "maximum nesting depth exceeded";
                if (message.entityId != null && Object.hasOwnProperty.call(message, "entityId"))
                    if (!$util.isString(message.entityId))
                        return "entityId: string expected";
                if (message["class"] != null && Object.hasOwnProperty.call(message, "class"))
                    switch (message["class"]) {
                    default:
                        return "class: enum value expected";
                    case 0:
                    case 1:
                    case 2:
                    case 3:
                    case 4:
                        break;
                    }
                if (message.latitude != null && Object.hasOwnProperty.call(message, "latitude"))
                    if (typeof message.latitude !== "number")
                        return "latitude: number expected";
                if (message.longitude != null && Object.hasOwnProperty.call(message, "longitude"))
                    if (typeof message.longitude !== "number")
                        return "longitude: number expected";
                if (message.altitudeMeters != null && Object.hasOwnProperty.call(message, "altitudeMeters"))
                    if (typeof message.altitudeMeters !== "number")
                        return "altitudeMeters: number expected";
                if (message.velocityKnots != null && Object.hasOwnProperty.call(message, "velocityKnots"))
                    if (typeof message.velocityKnots !== "number")
                        return "velocityKnots: number expected";
                if (message.headingDegrees != null && Object.hasOwnProperty.call(message, "headingDegrees"))
                    if (typeof message.headingDegrees !== "number")
                        return "headingDegrees: number expected";
                if (message.timestampMs != null && Object.hasOwnProperty.call(message, "timestampMs"))
                    if (!$util.isInteger(message.timestampMs) && !(message.timestampMs && $util.isInteger(message.timestampMs.low) && $util.isInteger(message.timestampMs.high)))
                        return "timestampMs: integer|Long expected";
                return null;
            };

            /**
             * Creates a Track message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {openlattice.telemetry.Track} Track
             */
            Track.fromObject = function fromObject(object, long) {
                if (object instanceof $root.openlattice.telemetry.Track)
                    return object;
                if (!$util.isObject(object))
                    throw TypeError(".openlattice.telemetry.Track: object expected");
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let message = new $root.openlattice.telemetry.Track();
                if (object.entityId != null)
                    message.entityId = String(object.entityId);
                switch (object["class"]) {
                default:
                    if (typeof object["class"] === "number") {
                        message["class"] = object["class"];
                        break;
                    }
                    break;
                case "UNKNOWN":
                case 0:
                    message["class"] = 0;
                    break;
                case "CIVILIAN_AIR":
                case 1:
                    message["class"] = 1;
                    break;
                case "CIVILIAN_MARITIME":
                case 2:
                    message["class"] = 2;
                    break;
                case "THREAT_AIR":
                case 3:
                    message["class"] = 3;
                    break;
                case "THREAT_MARITIME":
                case 4:
                    message["class"] = 4;
                    break;
                }
                if (object.latitude != null)
                    message.latitude = Number(object.latitude);
                if (object.longitude != null)
                    message.longitude = Number(object.longitude);
                if (object.altitudeMeters != null)
                    message.altitudeMeters = Number(object.altitudeMeters);
                if (object.velocityKnots != null)
                    message.velocityKnots = Number(object.velocityKnots);
                if (object.headingDegrees != null)
                    message.headingDegrees = Number(object.headingDegrees);
                if (object.timestampMs != null)
                    if ($util.Long)
                        message.timestampMs = $util.Long.fromValue(object.timestampMs, false);
                    else if (typeof object.timestampMs === "string")
                        message.timestampMs = parseInt(object.timestampMs, 10);
                    else if (typeof object.timestampMs === "number")
                        message.timestampMs = object.timestampMs;
                    else if (typeof object.timestampMs === "object")
                        message.timestampMs = new $util.LongBits(object.timestampMs.low >>> 0, object.timestampMs.high >>> 0).toNumber();
                return message;
            };

            /**
             * Creates a plain object from a Track message. Also converts values to other types if specified.
             * @function toObject
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {openlattice.telemetry.Track} message Track
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Track.toObject = function toObject(message, options, q) {
                if (!options)
                    options = {};
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                let object = {};
                if (options.defaults) {
                    object.entityId = "";
                    object["class"] = options.enums === String ? "UNKNOWN" : 0;
                    object.latitude = 0;
                    object.longitude = 0;
                    object.altitudeMeters = 0;
                    object.velocityKnots = 0;
                    object.headingDegrees = 0;
                    if ($util.Long) {
                        let long = new $util.Long(0, 0, false);
                        object.timestampMs = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : typeof BigInt !== "undefined" && options.longs === BigInt ? long.toBigInt() : long;
                    } else
                        object.timestampMs = options.longs === String ? "0" : typeof BigInt !== "undefined" && options.longs === BigInt ? BigInt("0") : 0;
                }
                if (message.entityId != null && Object.hasOwnProperty.call(message, "entityId"))
                    object.entityId = message.entityId;
                if (message["class"] != null && Object.hasOwnProperty.call(message, "class"))
                    object["class"] = options.enums === String ? $root.openlattice.telemetry.Track.TrackClass[message["class"]] === undefined ? message["class"] : $root.openlattice.telemetry.Track.TrackClass[message["class"]] : message["class"];
                if (message.latitude != null && Object.hasOwnProperty.call(message, "latitude"))
                    object.latitude = options.json && !isFinite(message.latitude) ? String(message.latitude) : message.latitude;
                if (message.longitude != null && Object.hasOwnProperty.call(message, "longitude"))
                    object.longitude = options.json && !isFinite(message.longitude) ? String(message.longitude) : message.longitude;
                if (message.altitudeMeters != null && Object.hasOwnProperty.call(message, "altitudeMeters"))
                    object.altitudeMeters = options.json && !isFinite(message.altitudeMeters) ? String(message.altitudeMeters) : message.altitudeMeters;
                if (message.velocityKnots != null && Object.hasOwnProperty.call(message, "velocityKnots"))
                    object.velocityKnots = options.json && !isFinite(message.velocityKnots) ? String(message.velocityKnots) : message.velocityKnots;
                if (message.headingDegrees != null && Object.hasOwnProperty.call(message, "headingDegrees"))
                    object.headingDegrees = options.json && !isFinite(message.headingDegrees) ? String(message.headingDegrees) : message.headingDegrees;
                if (message.timestampMs != null && Object.hasOwnProperty.call(message, "timestampMs"))
                    if (typeof BigInt !== "undefined" && options.longs === BigInt)
                        object.timestampMs = typeof message.timestampMs === "number" ? BigInt(message.timestampMs) : $util.Long.fromBits(message.timestampMs.low >>> 0, message.timestampMs.high >>> 0, false).toBigInt();
                    else if (typeof message.timestampMs === "number")
                        object.timestampMs = options.longs === String ? String(message.timestampMs) : message.timestampMs;
                    else
                        object.timestampMs = options.longs === String ? $util.Long.prototype.toString.call(message.timestampMs) : options.longs === Number ? new $util.LongBits(message.timestampMs.low >>> 0, message.timestampMs.high >>> 0).toNumber() : message.timestampMs;
                return object;
            };

            /**
             * Converts this Track to JSON.
             * @function toJSON
             * @memberof openlattice.telemetry.Track
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Track.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * Gets the default type url for Track
             * @function getTypeUrl
             * @memberof openlattice.telemetry.Track
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Track.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/openlattice.telemetry.Track";
            };

            /**
             * TrackClass enum.
             * @name openlattice.telemetry.Track.TrackClass
             * @enum {number}
             * @property {number} UNKNOWN=0 UNKNOWN value
             * @property {number} CIVILIAN_AIR=1 CIVILIAN_AIR value
             * @property {number} CIVILIAN_MARITIME=2 CIVILIAN_MARITIME value
             * @property {number} THREAT_AIR=3 THREAT_AIR value
             * @property {number} THREAT_MARITIME=4 THREAT_MARITIME value
             */
            Track.TrackClass = (function() {
                const valuesById = {}, values = Object.create(valuesById);
                values[valuesById[0] = "UNKNOWN"] = 0;
                values[valuesById[1] = "CIVILIAN_AIR"] = 1;
                values[valuesById[2] = "CIVILIAN_MARITIME"] = 2;
                values[valuesById[3] = "THREAT_AIR"] = 3;
                values[valuesById[4] = "THREAT_MARITIME"] = 4;
                return values;
            })();

            return Track;
        })();

        telemetry.TrackUpdateStream = (function() {

            /**
             * Properties of a TrackUpdateStream.
             * @memberof openlattice.telemetry
             * @interface ITrackUpdateStream
             * @property {Array.<openlattice.telemetry.ITrack>|null} [tracks] TrackUpdateStream tracks
             */

            /**
             * Constructs a new TrackUpdateStream.
             * @memberof openlattice.telemetry
             * @classdesc Represents a TrackUpdateStream.
             * @implements ITrackUpdateStream
             * @constructor
             * @param {openlattice.telemetry.ITrackUpdateStream=} [properties] Properties to set
             */
            function TrackUpdateStream(properties) {
                this.tracks = [];
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null && keys[i] !== "__proto__")
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * TrackUpdateStream tracks.
             * @member {Array.<openlattice.telemetry.ITrack>} tracks
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @instance
             */
            TrackUpdateStream.prototype.tracks = $util.emptyArray;

            /**
             * Creates a new TrackUpdateStream instance using the specified properties.
             * @function create
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {openlattice.telemetry.ITrackUpdateStream=} [properties] Properties to set
             * @returns {openlattice.telemetry.TrackUpdateStream} TrackUpdateStream instance
             */
            TrackUpdateStream.create = function create(properties) {
                return new TrackUpdateStream(properties);
            };

            /**
             * Encodes the specified TrackUpdateStream message. Does not implicitly {@link openlattice.telemetry.TrackUpdateStream.verify|verify} messages.
             * @function encode
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {openlattice.telemetry.ITrackUpdateStream} message TrackUpdateStream message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            TrackUpdateStream.encode = function encode(message, writer, q) {
                if (!writer)
                    writer = $Writer.create();
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                if (message.tracks != null && message.tracks.length)
                    for (let i = 0; i < message.tracks.length; ++i)
                        $root.openlattice.telemetry.Track.encode(message.tracks[i], writer.uint32(/* id 1, wireType 2 =*/10).fork(), q + 1).ldelim();
                return writer;
            };

            /**
             * Encodes the specified TrackUpdateStream message, length delimited. Does not implicitly {@link openlattice.telemetry.TrackUpdateStream.verify|verify} messages.
             * @function encodeDelimited
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {openlattice.telemetry.ITrackUpdateStream} message TrackUpdateStream message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            TrackUpdateStream.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer && writer.len ? writer.fork() : writer).ldelim();
            };

            /**
             * Decodes a TrackUpdateStream message from the specified reader or buffer.
             * @function decode
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {openlattice.telemetry.TrackUpdateStream} TrackUpdateStream
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            TrackUpdateStream.decode = function decode(reader, length, error, long) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                if (long === undefined)
                    long = 0;
                if (long > $Reader.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let end, message;
                if (length === undefined)
                    end = reader.len;
                else {
                    end = reader.pos + length;
                    if (end > reader.len)
                        throw RangeError("index out of range");
                    length = reader.len;
                    reader.len = end;
                }
                message = new $root.openlattice.telemetry.TrackUpdateStream();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    if (tag === error)
                        break;
                    switch (tag >>> 3) {
                    case 1: {
                            if (!(message.tracks && message.tracks.length))
                                message.tracks = [];
                            message.tracks.push($root.openlattice.telemetry.Track.decode(reader, reader.uint32(), undefined, long + 1));
                            break;
                        }
                    default:
                        reader.skipType(tag & 7, long);
                        break;
                    }
                }
                if (length !== undefined) {
                    if (reader.pos !== end)
                        throw RangeError("index out of range");
                    reader.len = length;
                }
                return message;
            };

            /**
             * Decodes a TrackUpdateStream message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {openlattice.telemetry.TrackUpdateStream} TrackUpdateStream
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            TrackUpdateStream.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a TrackUpdateStream message.
             * @function verify
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            TrackUpdateStream.verify = function verify(message, long) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    return "maximum nesting depth exceeded";
                if (message.tracks != null && Object.hasOwnProperty.call(message, "tracks")) {
                    if (!Array.isArray(message.tracks))
                        return "tracks: array expected";
                    for (let i = 0; i < message.tracks.length; ++i) {
                        let error = $root.openlattice.telemetry.Track.verify(message.tracks[i], long + 1);
                        if (error)
                            return "tracks." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a TrackUpdateStream message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {openlattice.telemetry.TrackUpdateStream} TrackUpdateStream
             */
            TrackUpdateStream.fromObject = function fromObject(object, long) {
                if (object instanceof $root.openlattice.telemetry.TrackUpdateStream)
                    return object;
                if (!$util.isObject(object))
                    throw TypeError(".openlattice.telemetry.TrackUpdateStream: object expected");
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let message = new $root.openlattice.telemetry.TrackUpdateStream();
                if (object.tracks) {
                    if (!Array.isArray(object.tracks))
                        throw TypeError(".openlattice.telemetry.TrackUpdateStream.tracks: array expected");
                    message.tracks = [];
                    for (let i = 0; i < object.tracks.length; ++i) {
                        if (!$util.isObject(object.tracks[i]))
                            throw TypeError(".openlattice.telemetry.TrackUpdateStream.tracks: object expected");
                        message.tracks[i] = $root.openlattice.telemetry.Track.fromObject(object.tracks[i], long + 1);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a TrackUpdateStream message. Also converts values to other types if specified.
             * @function toObject
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {openlattice.telemetry.TrackUpdateStream} message TrackUpdateStream
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            TrackUpdateStream.toObject = function toObject(message, options, q) {
                if (!options)
                    options = {};
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                let object = {};
                if (options.arrays || options.defaults)
                    object.tracks = [];
                if (message.tracks && message.tracks.length) {
                    object.tracks = [];
                    for (let j = 0; j < message.tracks.length; ++j)
                        object.tracks[j] = $root.openlattice.telemetry.Track.toObject(message.tracks[j], options, q + 1);
                }
                return object;
            };

            /**
             * Converts this TrackUpdateStream to JSON.
             * @function toJSON
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            TrackUpdateStream.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * Gets the default type url for TrackUpdateStream
             * @function getTypeUrl
             * @memberof openlattice.telemetry.TrackUpdateStream
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            TrackUpdateStream.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/openlattice.telemetry.TrackUpdateStream";
            };

            return TrackUpdateStream;
        })();

        telemetry.Ack = (function() {

            /**
             * Properties of an Ack.
             * @memberof openlattice.telemetry
             * @interface IAck
             * @property {boolean|null} [success] Ack success
             * @property {string|null} [message] Ack message
             */

            /**
             * Constructs a new Ack.
             * @memberof openlattice.telemetry
             * @classdesc Represents an Ack.
             * @implements IAck
             * @constructor
             * @param {openlattice.telemetry.IAck=} [properties] Properties to set
             */
            function Ack(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null && keys[i] !== "__proto__")
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Ack success.
             * @member {boolean} success
             * @memberof openlattice.telemetry.Ack
             * @instance
             */
            Ack.prototype.success = false;

            /**
             * Ack message.
             * @member {string} message
             * @memberof openlattice.telemetry.Ack
             * @instance
             */
            Ack.prototype.message = "";

            /**
             * Creates a new Ack instance using the specified properties.
             * @function create
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {openlattice.telemetry.IAck=} [properties] Properties to set
             * @returns {openlattice.telemetry.Ack} Ack instance
             */
            Ack.create = function create(properties) {
                return new Ack(properties);
            };

            /**
             * Encodes the specified Ack message. Does not implicitly {@link openlattice.telemetry.Ack.verify|verify} messages.
             * @function encode
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {openlattice.telemetry.IAck} message Ack message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Ack.encode = function encode(message, writer, q) {
                if (!writer)
                    writer = $Writer.create();
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                if (message.success != null && Object.hasOwnProperty.call(message, "success"))
                    writer.uint32(/* id 1, wireType 0 =*/8).bool(message.success);
                if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
                return writer;
            };

            /**
             * Encodes the specified Ack message, length delimited. Does not implicitly {@link openlattice.telemetry.Ack.verify|verify} messages.
             * @function encodeDelimited
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {openlattice.telemetry.IAck} message Ack message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Ack.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer && writer.len ? writer.fork() : writer).ldelim();
            };

            /**
             * Decodes an Ack message from the specified reader or buffer.
             * @function decode
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {openlattice.telemetry.Ack} Ack
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Ack.decode = function decode(reader, length, error, long) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                if (long === undefined)
                    long = 0;
                if (long > $Reader.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let end, message;
                if (length === undefined)
                    end = reader.len;
                else {
                    end = reader.pos + length;
                    if (end > reader.len)
                        throw RangeError("index out of range");
                    length = reader.len;
                    reader.len = end;
                }
                message = new $root.openlattice.telemetry.Ack();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    if (tag === error)
                        break;
                    switch (tag >>> 3) {
                    case 1: {
                            message.success = reader.bool();
                            break;
                        }
                    case 2: {
                            message.message = reader.string();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7, long);
                        break;
                    }
                }
                if (length !== undefined) {
                    if (reader.pos !== end)
                        throw RangeError("index out of range");
                    reader.len = length;
                }
                return message;
            };

            /**
             * Decodes an Ack message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {openlattice.telemetry.Ack} Ack
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Ack.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an Ack message.
             * @function verify
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Ack.verify = function verify(message, long) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    return "maximum nesting depth exceeded";
                if (message.success != null && Object.hasOwnProperty.call(message, "success"))
                    if (typeof message.success !== "boolean")
                        return "success: boolean expected";
                if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                    if (!$util.isString(message.message))
                        return "message: string expected";
                return null;
            };

            /**
             * Creates an Ack message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {openlattice.telemetry.Ack} Ack
             */
            Ack.fromObject = function fromObject(object, long) {
                if (object instanceof $root.openlattice.telemetry.Ack)
                    return object;
                if (!$util.isObject(object))
                    throw TypeError(".openlattice.telemetry.Ack: object expected");
                if (long === undefined)
                    long = 0;
                if (long > $util.recursionLimit)
                    throw Error("maximum nesting depth exceeded");
                let message = new $root.openlattice.telemetry.Ack();
                if (object.success != null)
                    message.success = Boolean(object.success);
                if (object.message != null)
                    message.message = String(object.message);
                return message;
            };

            /**
             * Creates a plain object from an Ack message. Also converts values to other types if specified.
             * @function toObject
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {openlattice.telemetry.Ack} message Ack
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Ack.toObject = function toObject(message, options, q) {
                if (!options)
                    options = {};
                if (q === undefined)
                    q = 0;
                if (q > $util.recursionLimit)
                    throw Error("max depth exceeded");
                let object = {};
                if (options.defaults) {
                    object.success = false;
                    object.message = "";
                }
                if (message.success != null && Object.hasOwnProperty.call(message, "success"))
                    object.success = message.success;
                if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                    object.message = message.message;
                return object;
            };

            /**
             * Converts this Ack to JSON.
             * @function toJSON
             * @memberof openlattice.telemetry.Ack
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Ack.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * Gets the default type url for Ack
             * @function getTypeUrl
             * @memberof openlattice.telemetry.Ack
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Ack.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/openlattice.telemetry.Ack";
            };

            return Ack;
        })();

        telemetry.ThreatService = (function() {

            /**
             * Constructs a new ThreatService service.
             * @memberof openlattice.telemetry
             * @classdesc Represents a ThreatService
             * @extends $protobuf.rpc.Service
             * @constructor
             * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
             * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
             * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
             */
            function ThreatService(rpcImpl, requestDelimited, responseDelimited) {
                $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
            }

            (ThreatService.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = ThreatService;

            /**
             * Creates new ThreatService service using the specified rpc implementation.
             * @function create
             * @memberof openlattice.telemetry.ThreatService
             * @static
             * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
             * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
             * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
             * @returns {ThreatService} RPC service. Useful where requests and/or responses are streamed.
             */
            ThreatService.create = function create(rpcImpl, requestDelimited, responseDelimited) {
                return new this(rpcImpl, requestDelimited, responseDelimited);
            };

            /**
             * Callback as used by {@link openlattice.telemetry.ThreatService#injectTrack}.
             * @memberof openlattice.telemetry.ThreatService
             * @typedef InjectTrackCallback
             * @type {function}
             * @param {Error|null} error Error, if any
             * @param {openlattice.telemetry.Ack} [response] Ack
             */

            /**
             * Calls InjectTrack.
             * @function injectTrack
             * @memberof openlattice.telemetry.ThreatService
             * @instance
             * @param {openlattice.telemetry.ITrack} request Track message or plain object
             * @param {openlattice.telemetry.ThreatService.InjectTrackCallback} callback Node-style callback called with the error, if any, and Ack
             * @returns {undefined}
             * @variation 1
             */
            Object.defineProperty(ThreatService.prototype.injectTrack = function injectTrack(request, callback) {
                return $protobuf.rpc.Service.prototype.rpcCall.call(this, injectTrack, $root.openlattice.telemetry.Track, $root.openlattice.telemetry.Ack, request, callback);
            }, "name", { value: "InjectTrack" });

            /**
             * Calls InjectTrack.
             * @function injectTrack
             * @memberof openlattice.telemetry.ThreatService
             * @instance
             * @param {openlattice.telemetry.ITrack} request Track message or plain object
             * @returns {Promise<openlattice.telemetry.Ack>} Promise
             * @variation 2
             */

            return ThreatService;
        })();

        return telemetry;
    })();

    return openlattice;
})();

export { $root as default };
