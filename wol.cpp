#include <arpa/inet.h>
#include <cstring>
#include <iostream>
#include <sstream>
#include <string>
#include <sys/socket.h>
#include <unistd.h>
#include <vector>

using namespace std;

const string VERSION = R"(v1.0.1)";
const string HELPMENU = R"(WakeOnLan by KopyTKG

    -h, --help                  	Display this help message
    -v, --version               	Display the version of wol

  * marks required switch

Any errors please report to: <https://github.com/kopytkg/wol/issues>

usage
$ wol xx:xx:xx:xx:xx:xx)";

vector<uint8_t> assembleMac(const vector<string> &macParts);
vector<uint8_t> assemble(const vector<uint8_t> &mac);

void wol(const string &sMac) {
  // Split the MAC address into parts
  vector<string> macParts;
  istringstream ss(sMac);
  string part;
  while (getline(ss, part, ':')) {
    macParts.push_back(part);
  }

  if (macParts.size() < 6) {
    cerr << "Invalid mac address" << endl;
    exit(1);
  }

  // Assemble the MAC address and message
  auto mac = assembleMac(macParts);
  auto msg = assemble(mac);

  // Set up the UDP address for broadcast
  sockaddr_in broadcastAddr;
  memset(&broadcastAddr, 0, sizeof(broadcastAddr));
  broadcastAddr.sin_family = AF_INET;
  broadcastAddr.sin_port = htons(40000);
  broadcastAddr.sin_addr.s_addr = inet_addr("255.255.255.255");

  // Create a UDP socket
  int sockfd = socket(AF_INET, SOCK_DGRAM, 0);
  if (sockfd < 0) {
    cerr << "Error while creating socket" << endl;
    exit(1);
  }

  // Set the socket option to allow broadcast
  int broadcastEnable = 1;
  if (setsockopt(sockfd, SOL_SOCKET, SO_BROADCAST, &broadcastEnable,
                 sizeof(broadcastEnable)) < 0) {
    cerr << "Error setting socket options: " << strerror(errno) << endl;
    close(sockfd);
    exit(1);
  }

  // Send the message
  ssize_t sentBytes =
      sendto(sockfd, msg.data(), msg.size(), 0,
             (struct sockaddr *)&broadcastAddr, sizeof(broadcastAddr));
  if (sentBytes < 0) {
    cerr << "Error sending message: " << strerror(errno) << endl;
    close(sockfd);
    exit(1);
  }

  cout << "Sending WoL to " << sMac << endl;

  // Close the socket
  close(sockfd);
}

vector<uint8_t> assembleMac(const vector<string> &macParts) {
  vector<uint8_t> mac(6);
  for (size_t i = 0; i < macParts.size(); ++i) {
    mac[i] = static_cast<uint8_t>(stoi(macParts[i], nullptr, 16));
  }
  return mac;
}

vector<uint8_t> assemble(const vector<uint8_t> &mac) {
  vector<uint8_t> msg(
      6 + 16 * mac.size()); // 6 bytes of magic + 16 repetitions of MAC
  fill(msg.begin(), msg.begin() + 6, 0xFF); // First 6 bytes are 0xFF
  for (size_t i = 0; i < 16; ++i) {
    copy(mac.begin(), mac.end(), msg.begin() + 6 + i * mac.size());
  }
  return msg;
}

// Example usage
int main(int argc, char *argv[]) {
  if (argc < 2) {
    cout << HELPMENU << endl;
    return 0;
  }
  auto item = string(argv[1]);
  if (item == "-v" || item == "--version") {
    cout << VERSION << endl;
    return 0;
  }
  if (item == "-h" || item == "--help") {
    cout << HELPMENU << endl;
    return 0;
  }

  wol(item);
  return 0;
}
