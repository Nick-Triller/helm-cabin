const statusIdToNameMap = {
    0: "Unknown",
    1: "Deployed",
    2: "Deleted",
    3: "Superseded",
    4: "Failed",
    5: "Deleting",
    6: "Pending install",
    7: "Pending upgrade",
    8: "Pending rollback",
};

export function statusIdToName(id) {
    return statusIdToNameMap[id];
}
