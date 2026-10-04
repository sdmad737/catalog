<h1 align="center">
  <br>
  <img src="assets/img/catalog.svg" width="160px">
  <br>
  Catalog
  <br>
</h1>
<p align="center" style="width: 100; margin-top: -30px;">
   <a href="https://github.com/sdmad737/catalog/tree/main/docs">Docs</a>
   |
   <a href="https://github.com/sdmad737/catalog">GitHub</a>
</p>



Catalog is a self-hosted inventory and organization system tailored to band equipment. It keeps instruments, stage gear, accessories, storage locations, photos, and maintenance records searchable without the overhead of an enterprise asset-management suite.

- _Simple_ - Catalog is designed to be simple and easy to use. No complicated setup or configuration required. Use either a single docker container, or deploy yourself by compiling the binary for your platform of choice.
- _Blazingly Fast_ - Catalog is written in Go, which makes it extremely fast and requires minimal resources to deploy. In general idle memory usage is less than 50MB for the whole container.
- _Portable_ - Catalog is designed to be portable and run on anywhere. We use SQLite and an embedded Web UI to make it easy to deploy, use, and backup.

## Project Status

Catalog is currently in early active development and is currently in **beta** stage. This means that the project may still be unstable and clunky. Overall, we are striving to not introduce any breaking changes and have checks in place to ensure migrations and upgrades are smooth. However, we do not guarantee that there will be no breaking changes. We will try to keep the documentation up to date as we make changes.

## Features

- Create and Manage _Items_ by providing a name and a description - That's it! Catalog requires only a few details to be provided to create an item, after that you can specify as much detail as you want, or hide away some of the things you won't ever need.
- Optional Details for Items include
    - Warranty Information
    - Sold To Information
    - Purchased From Information
    - Item Identifications (Serial, Model, etc)
    - Categorized Attachments (Images, Manuals, General)
    - Arbitrary/Custom Fields
- CSV Import/Export for quickly creating and managing items
- Custom Reporting
  - Bill of Materials Export
  - QR Code Label Generator
- Organize _Items_ by creating _Labels_ and _Locations_ and assigning them to items.
- Multi-Tenant Support - All users are placed in a group and can only see items in their group. Invite band members or share an instance with another team.


## Why Not Use Something Else?

There are many capable inventory systems, but band gear benefits from a lighter workflow: quick item entry, flexible locations, photos, QR labels, and maintenance history. Catalog focuses on those everyday tasks while keeping setup and operation straightforward.

### Spreadsheet

That's a fair point. If your needs can be fulfilled by a Spreadsheet, I'd suggest using that instead. I've found spreadsheets get pretty unwieldy when you have a lot of data, and it's hard to keep track of what's where. I also wanted to be able to search and filter my data in a more robust way than a spreadsheet can provide. I also wanted to leave the door open for more advanced features in the future like maintenance logs, moving label generators, and more.

### Snipe-It?

Snipe-IT is a strong choice for formal IT asset management. Catalog instead favors a smaller, more approachable workflow for instruments and production equipment, trading enterprise controls for quicker day-to-day use.
