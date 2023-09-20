/*
 *  ethernet_hal.h
 *
 *  Copyright 2013-2021 Michael Zillgith
 *
 *  This file is part of Platform Abstraction Layer (libpal)
 *  for libiec61850, libmms, and lib60870.
 */

#ifndef ETHERNET_HAL_H_
#define ETHERNET_HAL_H_

#include "hal_base.h"

#ifdef __cplusplus
extern "C" {
#endif

/*! \addtogroup hal
   *
   *  @{
   */

/**
 * @defgroup HAL_ETHERNET Direct access to the Ethernet layer (optional - required by GOOSE and Sampled Values)
 *
 * @{
 */


/**
 * \brief Opaque handle that represents an Ethernet "socket".
 */
typedef struct sEthernetSocket* EthernetSocket;

/** Opaque reference for a set of Ethernet socket handles */
typedef struct sEthernetHandleSet* EthernetHandleSet;

typedef enum {
    ETHERNET_SOCKET_MODE_PROMISC, /**<< receive all Ethernet messages */
    ETHERNET_SOCKET_MODE_ALL_MULTICAST, /**<< receive all multicast messages */
    ETHERNET_SOCKET_MODE_MULTICAST, /**<< receive only specific multicast messages */
    ETHERNET_SOCKET_MODE_HOST_ONLY /**<< receive only messages for the host */
} EthernetSocketMode;

/**
 * \brief Create a new connection handle set (EthernetHandleSet)
 *
 * \return new EthernetHandleSet instance
 */
PAL_API EthernetHandleSet
EthernetHandleSet_new(void);

/**
 * \brief add a socket to an existing handle set
 *
 * \param self the HandleSet instance
 * \param sock the socket to add
 */
PAL_API void
EthernetHandleSet_addSocket(EthernetHandleSet self, const EthernetSocket sock);

/**
 * \brief remove a socket from an existing handle set
 *
 * \param self the HandleSet instance
 * \param sock the socket to add
 */
PAL_API void
EthernetHandleSet_removeSocket(EthernetHandleSet self, const EthernetSocket sock);

/**
 * \brief wait for a socket to become ready
 *
 * This function is corresponding to the BSD socket select function.
 * The function will return after \p timeoutMs ms if no data is pending.
 *
 * \param self the HandleSet instance
 * \param timeoutMs in milliseconds (ms)
 * \return It returns the number of sockets on which data is pending
 *   or 0 if no data is pending on any of the monitored connections.
 *   The function shall return -1 if a socket error occures.
 */
PAL_API int
EthernetHandleSet_waitReady(EthernetHandleSet self, unsigned int timeoutMs);

/**
 * \brief destroy the EthernetHandleSet instance
 *
 * \param self the HandleSet instance to destroy
 */
PAL_API void
EthernetHandleSet_destroy(EthernetHandleSet self);

/**
 * \brief Return the MAC address of an Ethernet interface.
 *
 * The result are the six bytes that make up the Ethernet MAC address.
 *
 * \param interfaceId the ID of the Ethernet interface
 * \param addr pointer to a buffer to store the MAC address
 */
PAL_API void
Ethernet_getInterfaceMACAddress(const char* interfaceId, uint8_t* addr);

/**
 * \brief Create an Ethernet socket using the specified interface and
 * destination MAC address.
 *
 * \param interfaceId the ID of the Ethernet interface
 * \param destAddress byte array that contains the Ethernet destination MAC address for sending
 */
PAL_API EthernetSocket
Ethernet_createSocket(const char* interfaceId, uint8_t* destAddress);

/**
 * \brief destroy the ethernet socket
 *
 * \param ethSocket the ethernet socket handle
 */
PAL_API void
Ethernet_destroySocket(EthernetSocket ethSocket);

PAL_API void
Ethernet_sendPacket(EthernetSocket ethSocket, uint8_t* buffer, int packetSize);

/*
 * \brief set the receive mode of the Ethernet socket
 *
 * NOTE: When not implemented the the implementation has to be able to receive
 * all messages required by GOOSE and/or SV (usually multicast addresses).
 *
 * \param ethSocket the ethernet socket handle
 * \param mode the mode of the socket
 */
PAL_API void
Ethernet_setMode(EthernetSocket ethSocket, EthernetSocketMode mode);

/**
 * \brief Add a multicast address to be received by the Ethernet socket
 *
 * Used when mode is ETHERNET_SOCKET_MODE_MULTICAST
 *
 * \param ethSocket the ethernet socket handle
 * \param multicastAddress the multicast Ethernet address (this has to be a byte buffer of at least 6 byte)
 */
PAL_API void
Ethernet_addMulticastAddress(EthernetSocket ethSocket, uint8_t* multicastAddress);

/*
 * \brief set a protocol filter for the specified etherType
 *
 * NOTE: Implementation is not required but can improve the performance when the ethertype
 * filtering can be done on OS/network stack layer.
 *
 * \param ethSocket the ethernet socket handle
 * \param etherType the ether type of messages to accept
 */
PAL_API void
Ethernet_setProtocolFilter(EthernetSocket ethSocket, uint16_t etherType);

/**
 * \brief receive an ethernet packet (non-blocking)
 *
 * \param ethSocket the ethernet socket handle
 * \param buffer the buffer to copy the message to
 * \param bufferSize the maximum size of the buffer
 *
 * \return size of message received in bytes
 */
PAL_API int
Ethernet_receivePacket(EthernetSocket ethSocket, uint8_t* buffer, int bufferSize);

/**
 * \brief Indicates if runtime provides support for direct Ethernet access
 *
 * \return true if Ethernet support is available, false otherwise
 */
PAL_API bool
Ethernet_isSupported(void);

/*! @} */

/*! @} */

#ifdef __cplusplus
}
#endif

#endif /* ETHERNET_HAL_H_ */
