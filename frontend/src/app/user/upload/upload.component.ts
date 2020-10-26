import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { Subscription } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { FilesService, FileParams } from 'src/app/core/services/files';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class UploadComponent implements OnInit {

  public addEditFileForm: FormGroup;
  private subscriptions = new Subscription();
  public fileId: string;
  public showSpinner = false;
  constructor(
    private formBuilder: FormBuilder,
    private activatedRoute: ActivatedRoute,
    private snackBar: MatSnackBar,
    private router: Router,
    private filesService: FilesService
  ) {}

  ngOnInit() {
    this.activatedRoute.params.subscribe(data => {
      this.fileId = data.name;
    });
    this.addEditFileForm = this.formBuilder.group({
      name: [{value:'', disabled: true}, Validators.required, ],
      description: ['', Validators.required],
      file: ['', Validators.required],
    });
    if (this.fileId) {
      this.getFilesDetailsForEdit();
    }
}
openSnackBar(message: string, className: string ) {
  this.snackBar.open(message, '', {
    duration: 5000,
    panelClass: ['mat-toolbar', className]
  });
}
createPayload(): FileParams {
  const payload: FileParams = {
    filename: this.addEditFileForm.controls.name.value,
    file: this.addEditFileForm.controls.file.value,
    description: this.addEditFileForm.controls.description.value,
    UserID: sessionStorage.getItem('UserID')
  };
  return payload;
}

handleFileInput(files: FileList) {
  this.addEditFileForm.controls.file.setValue(files.item(0));
  this.addEditFileForm.controls.name.setValue(files.item(0).name);
}
  addEditFile() {
    this.showSpinner = true;
    const payload = this.createPayload();
    this.subscriptions.add(
      this.filesService.uploadFile(payload).subscribe(res => {
        this.openSnackBar('File uploaded','mat-accent');
        this.addEditFileForm.reset();
        this.showSpinner = false;
      }, error => {
        this.openSnackBar(error.error.error, 'mat-warn');
      }));
  }

  getFilesDetailsForEdit() {
    if (sessionStorage.getItem('file')) {
      const fileDetails = JSON.parse(sessionStorage.getItem('file'));
      this.addEditFileForm = this.formBuilder.group({
        name: [{value: fileDetails.file_name, disabled: true}, Validators.required],
        file: ['', Validators.required],
        description: [fileDetails.description, Validators.required]
      });
    }
  }

  editfile() {
    if(this.addEditFileForm.controls.file.value) {
      if(this.addEditFileForm.controls.name.value === this.fileId) {
    this.showSpinner = true;
    const payload = this.createPayload();
    this.subscriptions.add(
      this.filesService.editFile(payload).subscribe(res => {
        this.openSnackBar(res.data.message, 'mat-accent');
        this.showSpinner = false;
        this.router.navigate(['/home']);
      }, error => {
        this.openSnackBar(error.data.message, 'mat-warn');
      }));
    } else {
      this.openSnackBar('You cannot upload a new file. Please edit the same file and try again', 'mat-warn');
    }
    } else {
      this.openSnackBar('Please upload the updated file', 'mat-warn');
    }
  }


  ngOnDestroy() {
    this.subscriptions.unsubscribe();
    sessionStorage.removeItem('file');
  }
}
